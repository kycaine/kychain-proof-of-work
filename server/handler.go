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
	var tx blockchain.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction data"})
		return
	}

	bc.AddTransaction(tx.Sender, tx.Recipient, tx.Amount)
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
