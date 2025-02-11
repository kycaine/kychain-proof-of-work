package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	PreviousHash string
	Hash         string
	Nonce        int
}

func NewBlock(transactions []Transaction, previousHash string) *Block {
	block := &Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PreviousHash: previousHash,
	}
	block.Mine(getDifficulty())
	return block
}

func (b *Block) calculateHash() string {
	data := fmt.Sprintf("%d%s%s%d", b.Index, b.Timestamp, b.PreviousHash, b.Nonce)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
func getDifficulty() int {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	difficultyStr := os.Getenv("DIFFICULTY")
	difficulty, err := strconv.Atoi(difficultyStr)
	if err != nil {
		difficulty = 2
	}
	return difficulty
}
