package server

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	route := gin.Default()

	route.GET("/blockchain", GetBlockchainHandler)
	route.POST("/transaction", AddTransactionHandler)
	route.POST("/mine/:address", MineBlockHandler)
	route.GET("/balance/:address", GetBalanceHandler)

	return route
}
