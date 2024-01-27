// internal/api/handlers.go
package api

import (
	"Wallet/internal/database"
	"Wallet/internal/usecase"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetWalletStateHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		walletID := mux.Vars(r)["walletId"]
		if walletID == "" {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		balance, err := database.GetWalletBalance(db, walletID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"id":      walletID,
			"balance": balance,
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)
	}
}

func GetTransactionHistoryHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		walletID := mux.Vars(r)["walletId"]
		if walletID == "" {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		transactions, err := db.GetTransactionHistory(walletID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transactions)
	}
}

func CreateWalletHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		newWallet, err := usecase.CreateWallet(db)
		if err != nil {
			http.Error(w, "Internal Server Error CreateWallet", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newWallet)
	}
}

func SendMoneyHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(r)
		senderWalletID := vars["walletId"]
		if senderWalletID == "" {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		var sendRequest struct {
			To     string  `json:"to"`
			Amount float64 `json:"amount"`
		}

		err := json.NewDecoder(r.Body).Decode(&sendRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = usecase.TransferMoney(db, senderWalletID, sendRequest.To, sendRequest.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
