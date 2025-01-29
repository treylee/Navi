package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"navi-wings/internal/models"
	"navi-wings/config"
	"navi-wings/producer" // Import the producer package to access PublishMessage
	"github.com/gin-gonic/gin"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka" // Import confluent-kafka-go
)

func RegisterMessageRoutes(router *gin.Engine, kafkaProducer *kafka.Producer) {
	// Route to get all messages
	router.GET("/api/messages", GetMessages)

	// Route to post a message and publish it to Kafka
	router.POST("/api/messages", func(c *gin.Context) {
		// Handle message creation
		var newMessage models.Message
		
		// Bind JSON to message struct
		if err := c.ShouldBindJSON(&newMessage); err != nil {
			log.Printf("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Log the incoming message
		log.Printf("Received new message: %v\n", newMessage)

		// Save the message to the database
		if err := config.DB.Create(&newMessage).Error; err != nil {
			log.Printf("Error saving message to database: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
			return
		}

		// Log the successful database insertion
		log.Printf("Message saved to database: %v\n", newMessage)

		// After saving, publish the message to Kafka
		messageBytes, err := json.Marshal(newMessage)
		if err != nil {
			log.Printf("Error marshalling message: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize message"})
			return
		}

		// Log the Kafka message before publishing
		log.Printf("Publishing message to Kafka: %v\n", newMessage)

		// Publish the message to Kafka using the PublishMessage function from the producer package
		err = producer.PublishMessage(kafkaProducer, string(messageBytes))
		if err != nil {
			log.Printf("Error publishing message to Kafka: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message to Kafka"})
			return
		}

		// Log the successful Kafka publishing
		log.Printf("Message successfully published to Kafka: %v\n", newMessage)

		// Respond with the created message
		c.JSON(http.StatusOK, newMessage)
	})
}

// GetMessages retrieves all messages from the database and sends them as JSON
func GetMessages(c *gin.Context) {
	var messages []models.Message

	// Fetch messages from the database
	if err := config.DB.Find(&messages).Error; err != nil {
		log.Printf("Error fetching messages from database: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	// Return the messages as JSON
	c.JSON(http.StatusOK, messages)
}
