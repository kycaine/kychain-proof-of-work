package server

import (
	"net/http"
	"proof-of-work/blockchain"

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

	if len(BC.Mempool) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No transactions to mine"})
		return
	}

	newBlock := BC.MineBlock(address)
	BC.BroadcastBlock(*newBlock)

	c.JSON(http.StatusOK, gin.H{
		"message": "Block mined successfully",
		"block":   newBlock,
	})
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

	if !block.IsValidProof() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid proof-of-work"})
		return
	}

	BC.Blocks = append(BC.Blocks, &block)

	c.JSON(http.StatusOK, gin.H{"message": "Block received and added"})
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
