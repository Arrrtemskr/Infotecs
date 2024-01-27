// internal/usecase/wallet.go
package usecase

import (
	"Wallet/internal/database"
	"Wallet/internal/model"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func TransferMoney(db *database.Database, senderID, receiverID string, amount float64) error {

	senderBalance, err := database.GetWalletBalance(db, senderID)
	if err != nil {
		return err
	}

	receiverBalance, err := database.GetWalletBalance(db, receiverID)
	if err != nil {
		return err
	}

	if senderBalance < amount {
		return errors.New("insufficient funds")
	}

	senderBalance -= amount
	receiverBalance += amount

	err = database.UpdateWalletBalance(db, senderID, senderBalance)
	if err != nil {
		return err
	}

	err = database.UpdateWalletBalance(db, receiverID, receiverBalance)
	if err != nil {
		return err
	}

	err = database.LogTransaction(db, senderID, receiverID, amount)
	if err != nil {
		return err
	}

	return nil
}

func CreateWallet(db *database.Database) (*model.Wallet, error) {

	walletID, err := generateWalletID()
	if err != nil {
		return nil, err
	}

	newWallet := &model.Wallet{
		ID:      walletID,
		Balance: 100.0,
	}

	err = db.InsertWalletToDB(newWallet)
	if err != nil {
		return nil, err
	}

	return newWallet, nil
}

func generateWalletID() (string, error) {
	idBytes := make([]byte, 16)
	_, err := rand.Read(idBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(idBytes), nil
}
