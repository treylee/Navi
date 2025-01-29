package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors" // Import CORS middleware
	"navi-wings/config"
	"navi-wings/producer"  // Import producer to use CreateProducer
	"navi-wings/internal/routes"
	"time"
)

func main() {
	// Initialize configuration
	config.LoadConfig()

	// Initialize the database
	config.InitDatabase()

	// Initialize Kafka producer using confluent-kafka-go
	kafkaProducer, err := producer.CreateProducer() // Get Kafka producer from producer package
	if err != nil {
		log.Fatalf("Error initializing Kafka producer: %v", err)  // Fatal if Kafka producer can't be created
	}

	// Set up Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     viper.GetStringSlice("cors.allowedOrigins"),
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Register the routes and pass the Kafka producer to them
	routes.RegisterMessageRoutes(router, kafkaProducer)

	// Get the server port from config
	serverPort := viper.GetString("server.port")
	if serverPort == "" {
		log.Fatal("Server port not configured!")
	}

	// Set up graceful shutdown with signal handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start the Gin server in a separate goroutine
	go func() {
		log.Printf("Server started at http://localhost:%s", serverPort)
		if err := router.Run(":" + serverPort); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for a termination signal (SIGINT/SIGTERM)
	<-stop
	log.Println("Shutting down the server...")

	// Gracefully close Kafka producer
	if err := producer.CloseProducer(kafkaProducer); err != nil {
		log.Printf("Error closing Kafka producer: %v", err)
	} else {
		log.Println("Kafka producer closed gracefully.")
	}

	// Allow time for pending requests to finish before shutdown
	time.Sleep(2 * time.Second)

	log.Println("Server gracefully stopped")
}
