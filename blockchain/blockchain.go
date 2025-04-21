package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
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
	bc.CleanMempool(transactions)

	fmt.Println("Block mined successfully:", newBlock.Hash)

	bc.BroadcastBlock(*newBlock)

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
		jsonBlock, err := json.Marshal(block)
		if err != nil {
			fmt.Println("Error marshalling block:", err)
			continue
		}

		// asynchronus
		go func(nodeURL string, blockData []byte) {
			resp, err := http.Post(nodeURL, "application/json", bytes.NewBuffer(blockData))
			if err != nil {
				fmt.Println("Failed to broadcast to node:", nodeURL, "Error:", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				fmt.Println("Node rejected block:", nodeURL, "Status:", resp.Status)
			}
		}(url, jsonBlock)
	}

	fmt.Println("Block broadcast initiated to", len(bc.Nodes), "nodes")
}

func (bc *Blockchain) IsChainValid() bool {
	if len(bc.Blocks) == 1 {
		return true
	}

	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		previousBlock := bc.Blocks[i-1]

		if currentBlock.Hash != currentBlock.CalculateHash() {
			fmt.Println("Invalid hash for block", currentBlock.Index)
			return false
		}

		if currentBlock.PreviousHash != previousBlock.Hash {
			fmt.Println("Block", currentBlock.Index, "has invalid previous hash reference")
			return false
		}

		prefix := strings.Repeat("0", currentBlock.Difficulty)
		if !strings.HasPrefix(currentBlock.Hash, prefix) {
			fmt.Println("Block", currentBlock.Index, "has invalid proof of work")
			return false
		}

		if currentBlock.Index != previousBlock.Index+1 {
			fmt.Println("Block", currentBlock.Index, "has invalid index sequence")
			return false
		}
	}

	bc.ValidateTransactions()

	return true
}

func (bc *Blockchain) ValidateTransactions() bool {
	balances := make(map[string]float64)

	for _, tx := range bc.Blocks[0].Transactions {
		if tx.Recipient != "" {
			balances[tx.Recipient] += tx.Amount
		}
	}

	for i := 1; i < len(bc.Blocks); i++ {
		block := bc.Blocks[i]

		for _, tx := range block.Transactions {
			if tx.Sender == "mining-reward" {
				balances[tx.Recipient] += tx.Amount
				continue
			}

			if balances[tx.Sender] < (tx.Amount + tx.Fee) {
				fmt.Println("Invalid transaction in block", block.Index, "- insufficient funds")
				return false
			}

			balances[tx.Sender] -= (tx.Amount + tx.Fee)
			balances[tx.Recipient] += tx.Amount
		}
	}

	return true
}

func (bc *Blockchain) CleanMempool(blockTransactions []Transaction) {
	var newMempool []Transaction

	for _, mempoolTx := range bc.Mempool {
		found := false

		for _, blockTx := range blockTransactions {
			if mempoolTx.ID == blockTx.ID {
				found = true
				break
			}
		}

		if !found {
			newMempool = append(newMempool, mempoolTx)
		}
	}

	bc.Mempool = newMempool
}
