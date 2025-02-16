package database

import (
	"avitotech/internal/customErrors"
	"avitotech/internal/entities"
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		postgres.WithInitScripts("test_init.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(20*time.Second)),
	)

	if err != nil {
		return nil, err
	}

	database = dbName
	password = dbPwd
	username = dbUser

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	host = dbHost
	port = dbPort.Port()
	schema = "public"

	fmt.Println("host:", host, "port:", port)

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	fmt.Println("Started postgres container")
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
	fmt.Println("Stopped postgres container")
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestAddUser(t *testing.T) {
	srv := New()
	if err := srv.AddUser(&entities.User{Username: "user", Password: "password"}); err != nil {
		t.Fatalf("expected AddUser() to return nil, got %v", err)
	}
	if err := srv.AddUser(&entities.User{Username: "user", Password: "password"}); err == nil {
		t.Fatalf("expected AddUser() with duplicate username to return duplicate error, got %v", err)
	}

}

func TestGetUserByName(t *testing.T) {
	srv := New()
	user, err := srv.GetUserByName("unknownuser")

	if user != nil || err != nil {
		t.Fatalf("expected GetUserByName() with unknown user to return (nil, nil), got (%v, %v)", user, err)
	}

	err = srv.AddUser(&entities.User{Username: "knownuser", Password: "password"})
	if err != nil {
		t.Fatalf("Unexpected error while adding user: %v", err)
	}

	user, err = srv.GetUserByName("knownuser")

	if user == nil || err != nil {
		t.Fatalf("expected GetUserByName() with known user to return (user, nil), got (%v, %v)", user, err)
	}

}

func TestGetUserNameById(t *testing.T) {
	srv := New()
	err := srv.AddUser(&entities.User{Username: "knownuser", Password: "password"})
	if err != nil {
		t.Fatalf("Unexpected error while adding user: %v", err)
	}
	user, err := srv.GetUserByName("knownuser")
	if err != nil {
		t.Fatalf("Unexpected error while getting user: %v", err)
	}
	if username := srv.GetUserNameById(0); username != "<unknown>" {
		t.Fatalf("expected GetUserNameById() to return <unknown>, got %v", username)
	}
	if username = srv.GetUserNameById(user.ID); username != "knownuser" {
		t.Fatalf("expected GetUserNameById() to return knownuser, got %v", username)
	}
}

func TestGetCoinsByUserID(t *testing.T) {
	srv := New()
	err := srv.AddUser(&entities.User{Username: "knownuser", Password: "password"})
	if err != nil {
		t.Fatalf("Unexpected error while adding user: %v", err)
	}
	user, err := srv.GetUserByName("knownuser")
	if err != nil {
		t.Fatalf("Unexpected error while getting user: %v", err)
	}
	coins, err := srv.GetCoinsByUserID(user.ID)
	if err != nil {
		t.Fatalf("expected GetCoinsByUserID() don't return error, got %v", err)
	}
	if coins != INITIAL_COINS {
		t.Fatalf("expected GetCoinsByUserID() to return %v, got %v", INITIAL_COINS, coins)
	}
	coins, err = srv.GetCoinsByUserID(-1)
	if err != nil {
		t.Fatalf("expected GetCoinsByUserID() don't return error, got %v", err)
	}
	if coins != 0 {
		t.Fatalf("expected GetCoinsByUserID() to return %v, got %v", 0, coins)
	}
}

func TestGetInventoryByUserID(t *testing.T) {
	srv := New()
	err := srv.AddUser(&entities.User{Username: "knownuser", Password: "password"})
	if err != nil {
		t.Fatalf("Unexpected error while adding user: %v", err)
	}
	user, err := srv.GetUserByName("knownuser")
	if err != nil {
		t.Fatalf("Unexpected error while getting user: %v", err)
	}
	inventory, err := srv.GetInventoryByUserID(user.ID)
	if err != nil {
		t.Fatalf("expected GetInventoryByUserID() not return error, got %v", err)
	}
	if len(inventory) != 0 {
		t.Fatalf("expected GetInventoryByUserID() to return empty inventory, got %v", inventory)
	}
	err = srv.BuyItem(user.ID, "hoody")
	if err != nil {
		t.Fatal("Unexpected error while buying item")
	}
	inventory, err = srv.GetInventoryByUserID(user.ID)
	if err != nil {
		t.Fatalf("expected GetInventoryByUserID() not return error, got %v", err)
	}
	if len(inventory) != 1 {
		t.Fatalf("expected GetInventoryByUserID() to return 1 item, got %v", len(inventory))
	}
}

func TestGetTransactionsByUserID(t *testing.T) {
	srv := New()
	err := srv.AddUser(&entities.User{Username: "knownuser1", Password: "password"})
	if err != nil {
		t.Fatalf("Unexpected error while adding user: %v", err)
	}
	err = srv.AddUser(&entities.User{Username: "knownuser2", Password: "password"})
	if err != nil {
		t.Fatalf("Unexpected error while adding user: %v", err)
	}
	user1, err := srv.GetUserByName("knownuser1")
	if err != nil {
		t.Fatalf("Unexpected error while getting user: %v", err)
	}
	user2, err := srv.GetUserByName("knownuser2")
	if err != nil {
		t.Fatalf("Unexpected error while getting user: %v", err)
	}
	if transactions, err := srv.GetTransactionsByUserID(user1.ID); err != nil || len(transactions) != 0 {
		t.Fatalf("expected GetTransactionsByUserID() to return %v, got %v", nil, transactions)
	}
	if err := srv.SendCoin(user1.ID, user2.ID, 1200); !errors.Is(err, customErrors.ErrNotEnoughCoins) {
		t.Fatalf("expected SendCoin() to return ErrNotEnoughCoins, got %v", err)
	}
	if err := srv.SendCoin(user1.ID, user2.ID, 500); err != nil {
		t.Fatalf("expected SendCoin() to return nil, got %v", err)
	}
	if coins, err := srv.GetCoinsByUserID(user1.ID); err != nil || coins != 500 {
		t.Fatalf("expected GetCoinsByUserID() to return %v, got %v", 500, coins)
	}
	if coins, err := srv.GetCoinsByUserID(user2.ID); err != nil || coins != 1500 {
		t.Fatalf("expected GetCoinsByUserID() to return %v, got %v", 1500, coins)
	}
	transactions, err := srv.GetTransactionsByUserID(user1.ID)
	if err != nil {
		t.Fatalf("expected GetTransactionsByUserID() not return error, got %v", err)
	}
	if len(transactions) != 1 {
		t.Fatalf("expected GetTransactionsByUserID() to return %v transaction, got %v", 1, len(transactions))
	}
	if transactions[0].ToUserID != user2.ID {
		t.Fatalf("expected GetTransactionsByUserID() to return transaction to user2, got %v", transactions[0].ToUserID)
	}
	if transactions[0].Amount != 500 {
		t.Fatalf("expected GetTransactionsByUserID() to return transaction amount %v, got %v", 500, transactions[0].Amount)
	}
	transactions, err = srv.GetTransactionsByUserID(user2.ID)
	if err != nil {
		t.Fatalf("expected GetTransactionsByUserID() not return error, got %v", err)
	}
	if len(transactions) != 1 {
		t.Fatalf("expected GetTransactionsByUserID() to return %v transaction, got %v", 1, len(transactions))
	}
	if transactions[0].FromUserID != user1.ID {
		t.Fatalf("expected GetTransactionsByUserID() to return transaction to user1, got %v", transactions[0].FromUserID)
	}
	if transactions[0].Amount != 500 {
		t.Fatalf("expected GetTransactionsByUserID() to return transaction amount %v, got %v", 500, transactions[0].Amount)
	}
}

func TestBuyItem(t *testing.T) {
	srv := New()
	err := srv.AddUser(&entities.User{Username: "knownuser", Password: "password"})
	if err != nil {
		t.Fatalf("Unexpected error while adding user: %v", err)
	}
	user, err := srv.GetUserByName("knownuser")
	if err != nil {
		t.Fatalf("Unexpected error while getting user: %v", err)
	}
	if err := srv.BuyItem(user.ID, "hoody"); err != nil {
		t.Fatalf("expected BuyItem() to return nil, got %v", err)
	}
	inventory, err := srv.GetInventoryByUserID(user.ID)
	if err != nil {
		t.Fatalf("expected GetInventoryByUserID() not return error, got %v", err)
	}
	if len(inventory) != 1 && inventory[0].ItemType != "hoody" {
		t.Fatalf("expected GetInventoryByUserID() to return inventory with %v, got %v", "hoody", inventory)
	}
}

func TestClose(t *testing.T) {
	srv := New()

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
