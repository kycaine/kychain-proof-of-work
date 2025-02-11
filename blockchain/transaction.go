package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Transaction struct {
	ID        string
	Sender    string
	Recipient string
	Amount    float64
	Fee       float64
	Timestamp string
	Message   string
}

func NewTransaction(sender, recipient string, amount float64, message string) *Transaction {
	timestamp := time.Now().Format(time.RFC3339)
	fee := amount * 0.05
	data := fmt.Sprintf("%s%s%s", timestamp, sender, recipient)

	hash := sha256.Sum256([]byte(data))
	txID := hex.EncodeToString(hash[:])

	return &Transaction{
		ID:        txID,
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		Fee:       fee,
		Timestamp: timestamp,
		Message:   message,
	}
}
