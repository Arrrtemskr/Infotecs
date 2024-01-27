// cmd/main.go
package main

import (
	"Wallet/internal/api"
	"Wallet/internal/database"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

var (
	dbInstance *database.Database
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	cfgDB := database.ConfigDB{}

	err = env.Parse(&cfgDB)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfgDB.Host, cfgDB.Port, cfgDB.User, cfgDB.Password, cfgDB.Name)
	fmt.Println(connStr)

	dbInstance, err = database.ConnectDB(&cfgDB)
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}
	defer dbInstance.DB.Close()

	http.HandleFunc("/api/v1/wallet", api.CreateWalletHandler(dbInstance))

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/wallet/{walletId}/send", api.SendMoneyHandler(dbInstance)).Methods("POST")
	router.HandleFunc("/api/v1/wallet/{walletId}/history", api.GetTransactionHistoryHandler(dbInstance))
	router.HandleFunc("/api/v1/wallet/{walletId}", api.GetWalletStateHandler(dbInstance))

	http.Handle("/", router)

	fmt.Println("Server is listening on :8080")
	http.ListenAndServe(":8080", nil)
}
