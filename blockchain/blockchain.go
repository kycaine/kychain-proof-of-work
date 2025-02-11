package blockchain

import (
	"fmt"
	"time"
)

type Blockchain struct {
	Blocks  []*Block
	Mempool []Transaction
	Wallets map[string]*Wallet
}

func NewBlockchain() *Blockchain {
	timeNow := time.Now()
	initialDifficulty := 3

	genesisTransactions := []Transaction{
		*NewTransaction("genesis", "ky", 1000000, "Initial balance"),
	}

	genesisBlock := &Block{
		Index:        0,
		Timestamp:    timeNow.String(),
		Transactions: genesisTransactions,
		PreviousHash: "0",
		Difficulty:   initialDifficulty,
	}

	genesisBlock.Mine(initialDifficulty)

	return &Blockchain{
		Blocks:  []*Block{genesisBlock},
		Wallets: make(map[string]*Wallet),
	}
}

func (bc *Blockchain) AddTransaction(sender, recipient string, amount float64, message string) bool {
	balance := bc.GetBalance(sender)
	fee := amount * 0.05
	totalCost := amount + fee

	if balance < totalCost {
		fmt.Println("Transaksi gagal: saldo tidak mencukupi")
		return false
	}

	tx := NewTransaction(sender, recipient, amount, message)

	bc.Mempool = append(bc.Mempool, *tx)
	fmt.Println("Transaksi ditambahkan ke mempool!")
	return true
}

func (bc *Blockchain) MineBlock() *Block {
	if len(bc.Mempool) == 0 {
		fmt.Println("Mempool is empty, there are no transactions to mine.")
		return nil
	}

	previousBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(bc.Mempool, previousBlock)

	newBlock.Index = len(bc.Blocks)
	bc.Blocks = append(bc.Blocks, newBlock)
	bc.Mempool = []Transaction{}

	fmt.Println("Block mined successfully:", newBlock.Hash)
	return newBlock
}

func (bc *Blockchain) GetBalance(address string) float64 {
	balance := 0.0

	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			if tx.Sender == address {
				balance -= tx.Amount + tx.Fee
			}
			if tx.Recipient == address {
				balance += tx.Amount
			}
		}
	}
	return balance
}
