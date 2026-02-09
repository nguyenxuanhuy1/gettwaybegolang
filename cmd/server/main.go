package main

import (
	"gateway/config"
	"gateway/internal/router"
	"log"
)

func main() {
	// Initialize database
	config.ConnectDB()
	
	// Initialize Redis
	config.ConnectRedis()

	// Setup router with dependencies
	r := router.SetupRouter(config.DB, config.RedisClient)

	log.Println("Server is running on :8081")
	r.Run(":8081")
}
