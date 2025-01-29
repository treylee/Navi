package producer

import (
	"log"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"navi-wings/config"  // For accessing configuration (such as brokers, topic, etc.)
	"navi-wings/utils"   // Assuming utils is where your logger is defined
)

// CreateProducer initializes and returns a Kafka Producer instance
func CreateProducer() (*kafka.Producer, error) {
	// Get Kafka configuration from config
	kafkaTopic := config.GetKafkaTopic()  // Get Kafka topic from config
	logger := utils.GetLogger()           // Custom logger

	// Log the start of the producer creation using both log.Println and logger
	log.Println("Initializing Kafka producer...")
	logger.Println("Initializing Kafka producer...")

	producerConfig := kafka.ConfigMap{
		"bootstrap.servers":  config.GetKafkaBootstrapServers(), // Load Kafka server from config
		"sasl.username":      config.GetKafkaAPIKey(),           // API key
		"sasl.password":      config.GetKafkaAPISecret(),        // API secret
		"sasl.mechanism":     "PLAIN",                           // Kafka authentication mechanism
		"security.protocol":  "SASL_SSL",                        // Kafka security protocol
		"group.id":           "navi-dev-cg",                    // Unique consumer group name
		"auto.offset.reset":  "earliest",                         // Start reading from the earliest message
	}

	utils.LogInfo(fmt.Sprintf("Connecting to Kafka Cluster with the following details:"))
	utils.LogInfo(fmt.Sprintf("  Bootstrap Servers: %v", config.GetKafkaBootstrapServers()))
	utils.LogInfo(fmt.Sprintf("  SASL Username: %s", config.GetKafkaAPIKey())) // Don't log sensitive info in production
	utils.LogInfo(fmt.Sprintf("  SASL Mechanism: %s", "PLAIN"))
	utils.LogInfo(fmt.Sprintf("  Security Protocol: %s", "SASL_SSL"))
	utils.LogInfo(fmt.Sprintf("  Consumer Group ID: %s", "navi-ears-cg"))
	utils.LogInfo(fmt.Sprintf("  Offset Reset: %s", "earliest"))

	// Create the Kafka producer
	producer, err := kafka.NewProducer(&producerConfig)
	if err != nil {
		log.Printf("Error creating Kafka producer: %v", err) // Basic log output
		logger.Printf("Error creating Kafka producer: %v", err) // Detailed log output
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Log success
	log.Println("Kafka producer created successfully.")
	logger.Println("Kafka producer created successfully.")

	// Ensure kafkaTopic is being used
	log.Printf("Using Kafka topic: %s", kafkaTopic)
	logger.Printf("Using Kafka topic: %s", kafkaTopic)

	return producer, nil
}

// PublishMessage sends a message to Kafka
func PublishMessage(producer *kafka.Producer, message string) error {
	// Get the Kafka topic from config
	kafkaTopic := config.GetKafkaTopic()

	// Log the message being sent using both log.Println and logger
	log.Printf("Preparing to send message to Kafka topic: %s", kafkaTopic)
	logger := utils.GetLogger()
	logger.Printf("Preparing to send message to Kafka topic: %s", kafkaTopic)

	// Create a message to send to Kafka
	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaTopic,       // Pointer to the topic string
			Partition: kafka.PartitionAny, // Use the default partition
		},
		Value: []byte(message), // Message body (converted to byte slice)
	}

	// Produce message to Kafka (this sends the message asynchronously)
	err := producer.Produce(kafkaMessage, nil) // `nil` for no delivery report callback
	if err != nil {
		log.Printf("Error producing message to Kafka: %v", err) // Basic log output
		logger.Printf("Error producing message to Kafka: %v", err) // Detailed log output
		return err
	}

	// Log success
	log.Printf("Message successfully sent to Kafka topic: %s", kafkaTopic)
	logger.Printf("Message successfully sent to Kafka topic: %s", kafkaTopic)
	return nil
}

// CloseProducer gracefully shuts down the Kafka producer
func CloseProducer(producer *kafka.Producer) error {
	// Log the producer shutdown attempt using both log.Println and logger
	log.Println("Shutting down Kafka producer...")
	logger := utils.GetLogger()
	logger.Println("Shutting down Kafka producer...")

	// Close the producer and flush any remaining messages
	producer.Close()

	// Log successful shutdown
	log.Println("Kafka producer closed gracefully.")
	logger.Println("Kafka producer closed gracefully.")
	return nil
}
