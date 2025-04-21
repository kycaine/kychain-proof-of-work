package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	PreviousHash string
	Hash         string
	Nonce        int
	Difficulty   int
}

func NewBlock(transactions []Transaction, previousBlock *Block) *Block {
	startTime := time.Now()

	difficulty := AdjustDifficulty(*previousBlock, startTime)

	block := &Block{
		Index:        previousBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PreviousHash: previousBlock.Hash,
		Difficulty:   difficulty,
	}

	block.Mine(difficulty)
	return block
}

// func (b *Block) calculateHash() string {
// 	data := fmt.Sprintf("%d%s%s%d", b.Index, b.Timestamp, b.PreviousHash, b.Nonce)
// 	hash := sha256.Sum256([]byte(data))
// 	return hex.EncodeToString(hash[:])
// }

func (b *Block) CalculateHash() string {
	var txDataBuffer bytes.Buffer

	for _, tx := range b.Transactions {
		txData := fmt.Sprintf("%s%s%s%.8f%.8f%s%s",
			tx.ID, tx.Sender, tx.Recipient, tx.Amount, tx.Fee, tx.Timestamp, tx.Message)
		txDataBuffer.WriteString(txData)
	}

	blockData := fmt.Sprintf("%d%s%s%s%d%d",
		b.Index, b.Timestamp, txDataBuffer.String(), b.PreviousHash, b.Nonce, b.Difficulty)

	hash := sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hash[:])
}

func AdjustDifficulty(previousBlock Block, startTime time.Time) int {
	targetTime := 300
	elapsedTime := int(time.Since(startTime).Seconds())

	if elapsedTime < targetTime {
		return previousBlock.Difficulty + 1
	} else if elapsedTime > targetTime {
		return previousBlock.Difficulty - 1
	}
	return previousBlock.Difficulty
}

func (b *Block) IsValidProof() bool {
	hash := b.CalculateHash()
	return strings.HasPrefix(hash, "0000")
}
