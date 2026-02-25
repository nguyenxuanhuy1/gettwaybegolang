package main

import (
	"gateway/config"
	"gateway/internal/router"
	"log"
	"os"
)

func main() {
	config.ConnectDB()
	r := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server is running on :" + port)
	r.Run(":" + port)
}
