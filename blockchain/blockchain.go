package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

type Blockchain struct {
	Blocks  []*Block
	Mempool []Transaction
	Wallets map[string]*Wallet
	Nodes   []string
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
		fmt.Println("Transaction failed: not enough balance")
		return false
	}

	tx := NewTransaction(sender, recipient, amount, message)

	bc.Mempool = append(bc.Mempool, *tx)
	fmt.Println("Transaction added to mempool!")
	return true
}

func (bc *Blockchain) MineBlock(minerAddress string) *Block {
	previousBlock := bc.Blocks[len(bc.Blocks)-1]

	rewardTx := Transaction{
		ID:        "reward-" + fmt.Sprint(len(bc.Blocks)),
		Sender:    "mining-reward",
		Recipient: minerAddress,
		Amount:    GetBlockReward(len(bc.Blocks)),
		Fee:       0,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   "Block reward",
	}

	transactions := append([]Transaction{rewardTx}, bc.Mempool...)

	newBlock := NewBlock(transactions, previousBlock)
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

func GetBlockReward(blockHeight int) float64 {
	baseReward := 100.0
	halvingInterval := 1000

	halvings := blockHeight / halvingInterval
	reward := baseReward / math.Pow(2, float64(halvings))

	return reward
}

func (bc *Blockchain) RegisterNode(address string) {
	for _, node := range bc.Nodes {
		if node == address {
			fmt.Println("Node already exists!")
			return
		}
	}
	bc.Nodes = append(bc.Nodes, address)
	fmt.Println("Node registered:", address)
}

func (bc *Blockchain) BroadcastBlock(block Block) {
	for _, node := range bc.Nodes {
		url := node + "/receive_block"
		jsonBlock, _ := json.Marshal(block)
		http.Post(url, "application/json", bytes.NewBuffer(jsonBlock))
	}
}
