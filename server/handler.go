package server

import (
	"net/http"
	"proof-of-work/blockchain"

	"github.com/gin-gonic/gin"
)

var bc = blockchain.NewBlockchain()

func GetBlockchainHandler(c *gin.Context) {
	c.JSON(http.StatusOK, bc.Blocks)
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

	bc.AddTransaction(request.Sender, request.Recipient, request.Amount, request.Message)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction added to mempool!"})
}

func MineBlockHandler(c *gin.Context) {
	newBlock := bc.MineBlock()
	if newBlock == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No transactions to mine!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Block mined successfully", "block": newBlock})
}

func GetBalanceHandler(c *gin.Context) {
	address := c.Param("address")

	balance := bc.GetBalance(address)
	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"balance": balance,
	})
}
