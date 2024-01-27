// internal/database/db.go
package database

import (
	"Wallet/internal/model"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type ConfigDB struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Name     string `env:"DB_NAME,required"`
}

type Database struct {
	DB *sql.DB
}

func ConnectDB(cfgDB *ConfigDB) (*Database, error) {

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfgDB.Host, cfgDB.Port, cfgDB.User, cfgDB.Password, "postgres")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Failed to connect to PostgreSQL:", err)
		return nil, err
	}

	if err := createTables(db); err != nil {
		log.Println("Failed to create tables:", err)
		return nil, err
	}

	fmt.Println("Connected to the database")

	return &Database{DB: db}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Wallets (
			ID VARCHAR(50) PRIMARY KEY,
			amount DECIMAL(10, 2)
		);

		CREATE TABLE IF NOT EXISTS History (
			id SERIAL PRIMARY KEY,
			time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			from_wallet_id VARCHAR(50),
			to_wallet_id VARCHAR(50),
			amount DECIMAL(10, 2)
		);
	`)
	return err
}

func (db *Database) InsertWalletToDB(wallet *model.Wallet) error {
	query := "INSERT INTO Wallets (ID, Amount) VALUES ($1, $2)"
	_, err := db.DB.Exec(query, wallet.ID, wallet.Balance)
	if err != nil {
		return fmt.Errorf("failed to insert wallet into DB: %v", err)
	}
	return nil
}

func GetWalletBalance(db *Database, walletID string) (float64, error) {
	var balance float64
	query := "SELECT Amount FROM Wallets WHERE ID = $1"
	err := db.DB.QueryRow(query, walletID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to get wallet balance: %w", err)
	}
	return balance, nil
}

func UpdateWalletBalance(db *Database, walletID string, newBalance float64) error {
	query := "UPDATE Wallets SET Amount = $1 WHERE ID = $2"
	_, err := db.DB.Exec(query, newBalance, walletID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}
	return nil
}

func LogTransaction(db *Database, from, to string, amount float64) error {
	query := "INSERT INTO History (time, from_wallet_id, to_wallet_id, amount) VALUES ($1, $2, $3, $4)"
	_, err := db.DB.Exec(query, time.Now(), from, to, amount)
	if err != nil {
		return fmt.Errorf("failed to log transaction: %w", err)
	}
	return nil
}

type Transaction struct {
	Time   time.Time `json:"time"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount float64   `json:"amount"`
}

func (db *Database) GetTransactionHistory(walletID string) ([]Transaction, error) {
	query := `
		SELECT time, from_wallet_id, to_wallet_id, amount
		FROM History
		WHERE from_wallet_id = $1 OR to_wallet_id = $1
	`

	rows, err := db.DB.Query(query, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction history: %w", err)
	}

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.Time, &transaction.From, &transaction.To, &transaction.Amount); err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction history rows: %w", err)
	}

	return transactions, nil
}
