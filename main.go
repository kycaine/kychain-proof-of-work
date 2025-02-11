package main

import (
	"log"
	"os"
	"proof-of-work/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := server.SetupRouter()
	r.Run(":" + port)
}
