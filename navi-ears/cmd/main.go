package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"navi-ears/config"
	"navi-ears/consumer"
	"navi-ears/utils" // Custom logger
	"github.com/confluentinc/confluent-kafka-go/v2/kafka" // Kafka package (v2)
	"time"
)

func main() {
	// Initialize logger
	logInstance := utils.GetLogger()

	// Load configuration
	config.LoadConfig()

	// Create Kafka consumer
	consumer, err := consumer.CreateConsumer()
	if err != nil {
		logInstance.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consume	r.Close()

	// Get the Kafka topic from config
	topic := config.GetKafkaTopic()

	// Subscribe to the topic (using the correct Kafka consumer method)
	err = consumer.Subscribe(topic, nil)  // Standard Kafka method
	if err != nil {
		logInstance.Fatalf("Error subscribing to Kafka topic: %v", err)
	}
	logInstance.Printf("Subscribed to Kafka topic: %s", topic)

	// Start consuming messages (using the standard Kafka consumer method)
	go consumeMessages(consumer, logInstance)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Set a timeout for graceful shutdown
	shutdownTimeout := time.After(30 * time.Second) // Timeout after 30 seconds if no signal

	select {
	case <-stop:
		logInstance.Println("Received shutdown signal. Shutting down gracefully...")
	case <-shutdownTimeout:
		logInstance.Println("Shutdown timeout reached. Forcefully shutting down...")
	}

	// Final cleanup
	logInstance.Println("Consumer shutdown completed.")
}

// consumeMessages consumes messages from Kafka and handles them
func consumeMessages(consumer *kafka.Consumer, logInstance *log.Logger) {
	for {
		// Consume message
		msg, err := consumer.ReadMessage(-1) // Blocking call, waits for messages
		if err != nil {
			logInstance.Printf("Error consuming message: %v", err)
			continue
		}

		// Handle consumed message
		logInstance.Printf("Consumed message: %s", string(msg.Value))
	}
}
