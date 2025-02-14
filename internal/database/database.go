package database

import (
	"avitotech/internal/entities"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"log/slog"
	"os"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
	// GetUserByName retrieves the user by the given username.
	GetUserByName(username string) (*entities.User, error)
	// AddUser inserts a new user into the database.
	AddUser(user *entities.User) error
}

type service struct {
	db *sql.DB
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
	dbInstance = &service{
		db: db,
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
	user := &entities.User{}
	row := s.db.QueryRow("SELECT id, username, password, created_at, updated_at FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// AddUser inserts a new user into the database.
func (s *service) AddUser(user *entities.User) error {
	_, err := s.db.Exec("INSERT INTO users (username, password, created_at, updated_at) VALUES ($1, $2, $3, $4)", user.Username, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
