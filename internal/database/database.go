package database

import (
	"avitotech/internal/customErrors"
	"avitotech/internal/entities"
	"avitotech/pkg/imcache"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
	// GetUserByName retrieves the user by the given username.
	GetUserByName(username string) (*entities.User, error)
	// GetUserNameById retrieves the username by the given user ID.
	GetUserNameById(userId int) string
	// AddUser inserts a new user into the database.
	AddUser(user *entities.User) error
	// GetCoinsByUserID retrieves the number of coins by the given user ID
	GetCoinsByUserID(userId int) (int, error)
	// GetInventoryByUserID retrieves the inventory items by the given user ID.
	GetInventoryByUserID(userId int) ([]entities.InventoryItem, error)
	// GetTransactionsByUserID retrieves the transactions by the given user ID.
	GetTransactionsByUserID(userId int) ([]entities.Transaction, error)
	// SendCoin sends coins from one user to another.
	SendCoin(fromUserID, toUserID, amount int) error
	// BuyItem buys an item for the given user.
	BuyItem(userId int, itemType string) error
}

type service struct {
	db    *sql.DB
	cache imcache.Cache
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

const (
	INITIAL_COINS = 1000
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	cache := imcache.NewInMemoryCache(5 * time.Minute)
	dbInstance = &service{
		db:    db,
		cache: cache,
	}
	return dbInstance
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	slog.Info("Disconnected from", "database", database)
	return s.db.Close()
}

// GetUserByName retrieves the user by the given username.
func (s *service) GetUserByName(username string) (*entities.User, error) {
	if user, ok := s.cache.Get(username); ok {
		return user.(*entities.User), nil
	}
	user := &entities.User{}
	row := s.db.QueryRow("SELECT id, username, password, created_at, updated_at FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	s.cache.Set(username, user)
	return user, nil
}

// GetUserNameById retrieves the username by the given user ID.
func (s *service) GetUserNameById(userId int) string {
	var username string
	row := s.db.QueryRow("SELECT username FROM users WHERE id = $1", userId)
	err := row.Scan(&username)
	if err != nil {
		return "<unknown>"
	}
	return username
}

// AddUser inserts a new user into the database.
func (s *service) AddUser(user *entities.User) error {
	var userId int
	err := s.db.QueryRow("INSERT INTO users (username, password, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", user.Username, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&userId)
	if err != nil {
		return err
	}
	if err := s.InitUserWallet(userId); err != nil {
		return err
	}
	slog.Info("User created and added coins to his wallet", "username", user.Username, "coins", INITIAL_COINS)
	return nil
}

// GetCoinsByUserID retrieves the number of coins by the given user ID.
func (s *service) GetCoinsByUserID(userId int) (int, error) {
	var coins int
	row := s.db.QueryRow("SELECT amount FROM coins WHERE user_id = $1", userId)
	err := row.Scan(&coins)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return coins, nil
}

// GetInventoryByUserID retrieves the inventory items by the given user ID.
func (s *service) GetInventoryByUserID(userId int) ([]entities.InventoryItem, error) {
	if inventoryItems, ok := s.cache.Get(strconv.Itoa(userId)); ok {
		return inventoryItems.([]entities.InventoryItem), nil
	}
	var inventoryItems []entities.InventoryItem
	rows, err := s.db.Query("SELECT item_type, quantity FROM inventory WHERE user_id = $1", userId)
	if errors.Is(err, sql.ErrNoRows) {
		return inventoryItems, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item entities.InventoryItem
		err = rows.Scan(&item.ItemType, &item.Quantity)

		if err != nil {
			return nil, err
		}
		inventoryItems = append(inventoryItems, item)
	}
	s.cache.Set(strconv.Itoa(userId), inventoryItems)
	return inventoryItems, nil
}

// GetTransactionsByUserID retrieves the transactions by the given user ID.
func (s *service) GetTransactionsByUserID(userId int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	rows, err := s.db.Query("SELECT from_user_id, to_user_id, amount FROM coin_transactions WHERE from_user_id = $1 OR to_user_id = $1", userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var transaction entities.Transaction
		err = rows.Scan(&transaction.FromUserID, &transaction.ToUserID, &transaction.Amount)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// SendCoin sends coins from one user to another.
func (s *service) SendCoin(fromUserID, toUserID, amount int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	fromUserCoins, err := s.GetCoinsByUserID(fromUserID)
	if err != nil {
		return err
	}

	if fromUserCoins < amount {
		return customErrors.ErrNotEnoughCoins
	}

	toUserCoins, err := s.GetCoinsByUserID(toUserID)
	if err != nil {
		return err
	}

	if err := s.UpdateCoins(fromUserID, fromUserCoins-amount); err != nil {
		return err
	}
	if err := s.UpdateCoins(toUserID, toUserCoins+amount); err != nil {
		return err
	}

	transaction := &entities.Transaction{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
	}
	if err := s.SaveTransaction(transaction); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdateCoins updates the number of coins by the given user ID.
func (s *service) UpdateCoins(userId, coins int) error {
	_, err := s.db.Exec("UPDATE coins SET amount = $1 WHERE user_id = $2", coins, userId)
	if err != nil {
		return err
	}
	return nil
}

// SaveTransaction inserts a new transaction into the database.
func (s *service) SaveTransaction(transaction *entities.Transaction) error {
	_, err := s.db.Exec("INSERT INTO coin_transactions (from_user_id, to_user_id, amount, transaction_type, created_at) VALUES ($1, $2, $3, $4, $5)", transaction.FromUserID, transaction.ToUserID, transaction.Amount, "send", time.Now())
	if err != nil {
		return err
	}
	return nil
}

// BuyItem buys an item for the given user.
func (s *service) BuyItem(userId int, itemType string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	coins, err := s.GetCoinsByUserID(userId)
	if err != nil {
		return err
	}

	itemPrice, err := s.GetItemPrice(itemType)
	if err != nil {
		return err
	}

	if coins < itemPrice {
		return customErrors.ErrNotEnoughCoins
	}

	if err := s.UpdateCoins(userId, coins-itemPrice); err != nil {
		return err
	}

	if err := s.AddItemToInventory(userId, itemType); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetItemPrice retrieves the price of the item by the given item type.
func (s *service) GetItemPrice(itemType string) (int, error) {
	var price int
	row := s.db.QueryRow("SELECT price FROM shop WHERE item_type = $1", itemType)
	err := row.Scan(&price)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, customErrors.ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	return price, nil
}

// AddItemToInventory adds an item to the inventory of the given user.
func (s *service) AddItemToInventory(userId int, itemType string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var quantity int
	err = s.db.QueryRow("SELECT quantity FROM inventory WHERE user_id = $1 AND item_type = $2", userId, itemType).Scan(&quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = s.db.Exec("INSERT INTO inventory (user_id, item_type, quantity) VALUES ($1, $2, 1)", userId, itemType)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		_, err = s.db.Exec("UPDATE inventory SET quantity = quantity + 1 WHERE user_id = $1 AND item_type = $2", userId, itemType)
		if err != nil {
			return err
		}
	}
	s.cache.Delete(strconv.Itoa(userId))
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// InitUserWallet initializes the user wallet with the initial amount of coins.
func (s *service) InitUserWallet(userId int) error {
	if _, err := s.db.Exec("INSERT INTO coins (user_id, amount) VALUES ($1, $2)", userId, INITIAL_COINS); err != nil {
		return err
	}
	return nil
}
