package server

import (
	"net/http"
	"proof-of-work/blockchain"
	"strings"

	"github.com/gin-gonic/gin"
)

var BC = blockchain.NewBlockchain()

func GetBlockchainHandler(c *gin.Context) {
	c.JSON(http.StatusOK, BC.Blocks)
}

func AddTransactionHandler(c *gin.Context) {
	var request struct {
		Sender    string  `json:"sender"`
		Recipient string  `json:"recipient"`
		Amount    float64 `json:"amount"`
		Message   string  `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction data"})
		return
	}

	BC.AddTransaction(request.Sender, request.Recipient, request.Amount, request.Message)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction added to mempool!"})
}

func MineBlockHandler(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address required"})
		return
	}

	block := BC.MineBlock(address)
	if block == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mine block"})
		return
	}

	c.JSON(http.StatusOK, block)
}

func GetBalanceHandler(c *gin.Context) {
	address := c.Param("address")
	balance := BC.GetBalance(address)

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"balance": balance,
	})
}

func ReceiveBlock(c *gin.Context) {
	var block blockchain.Block
	if err := c.ShouldBindJSON(&block); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block data"})
		return
	}

	lastBlock := BC.Blocks[len(BC.Blocks)-1]

	if block.Index <= lastBlock.Index {
		c.JSON(http.StatusConflict, gin.H{"error": "Block rejected: invalid index"})
		return
	}

	if block.PreviousHash != lastBlock.Hash {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid previous hash"})
		return
	}

	prefix := strings.Repeat("0", block.Difficulty)
	if !strings.HasPrefix(block.Hash, prefix) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid proof-of-work"})
		return
	}

	calculatedHash := block.CalculateHash()
	if calculatedHash != block.Hash {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid block hash"})
		return
	}

	if !validateBlockTransactions(&block, BC) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid transactions in block"})
		return
	}

	BC.Blocks = append(BC.Blocks, &block)

	if !BC.IsChainValid() {
		BC.Blocks = BC.Blocks[:len(BC.Blocks)-1]
		c.JSON(http.StatusBadRequest, gin.H{"error": "Block rejected: causes invalid blockchain state"})
		return
	}

	BC.CleanMempool(block.Transactions)

	c.JSON(http.StatusOK, gin.H{
		"message": "Block received and added successfully",
		"block":   block,
	})
}

func validateBlockTransactions(block *blockchain.Block, bc *blockchain.Blockchain) bool {
	startIdx := 0
	if len(block.Transactions) > 0 && block.Transactions[0].Sender == "mining-reward" {
		startIdx = 1
	}

	for i := startIdx; i < len(block.Transactions); i++ {
		tx := block.Transactions[i]

		if tx.Sender == "genesis" || tx.Sender == "mining-reward" {
			continue
		}

		senderBalance := bc.GetBalance(tx.Sender)
		if senderBalance < (tx.Amount + tx.Fee) {
			return false
		}

	}

	return true
}

func RegisterNodeHandler(c *gin.Context) {
	var request struct {
		Address string `json:"address"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	BC.RegisterNode(request.Address)
	c.JSON(http.StatusOK, gin.H{"message": "Node registered", "node": request.Address})
}
