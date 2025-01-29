package consumer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"navi-ears/utils"   // Assuming utils is where your logger is defined
	"navi-ears/config"   // Assuming config is set up via viper or environment variables
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Create and configure Kafka consumer from viper configuration
func CreateConsumer() (*kafka.Consumer, error) {
	// Load Kafka configuration from config
	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers":  config.GetKafkaBootstrapServers(), // Load Kafka server from config
		"sasl.username":      config.GetKafkaAPIKey(),           // API key
		"sasl.password":      config.GetKafkaAPISecret(),        // API secret
		"sasl.mechanism":     "PLAIN",                           // Kafka authentication mechanism
		"security.protocol":  "SASL_SSL",                        // Kafka security protocol
		"group.id":           "navi-dev-cg",                   // Unique consumer group name
		"auto.offset.reset":  "earliest",                        // Start reading from the earliest message
	}

	// Log the Kafka cluster connection details for debugging
	utils.LogInfo(fmt.Sprintf("Connecting to Kafka Cluster with the following details:"))
	utils.LogInfo(fmt.Sprintf("  Bootstrap Servers: %v", config.GetKafkaBootstrapServers()))
	utils.LogInfo(fmt.Sprintf("  SASL Username: %s", config.GetKafkaAPIKey())) // Don't log sensitive info in production
	utils.LogInfo(fmt.Sprintf("  SASL Mechanism: %s", "PLAIN"))
	utils.LogInfo(fmt.Sprintf("  Security Protocol: %s", "SASL_SSL"))
	utils.LogInfo(fmt.Sprintf("  Consumer Group ID: %s", "navi-ears-cg"))
	utils.LogInfo(fmt.Sprintf("  Offset Reset: %s", "earliest"))

	// Initialize the Kafka consumer
	consumer, err := kafka.NewConsumer(&kafkaConfig)
	if err != nil {
		utils.LogError("Failed to create Kafka consumer", err)  // Logging with utils package
		return nil, err
	}
	utils.LogInfo("Kafka consumer created successfully")
	return consumer, nil
}

// SubscribeToTopic subscribes to a Kafka topic and starts consuming messages
func SubscribeToTopic(consumer *kafka.Consumer, topic string) error {
	err := consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		utils.LogError("Failed to subscribe to Kafka topic", err)  // Logging with utils package
		return err
	}
	utils.LogInfo(fmt.Sprintf("Subscribed to Kafka topic: %s", topic))
	return nil
}

// ConsumeMessages reads messages from a Kafka topic and processes them
func ConsumeMessages(consumer *kafka.Consumer, topic string) {
	// Log that the consumer has started
	utils.LogInfo(fmt.Sprintf("Consumer started and waiting for messages from topic %s...", topic)) // Logging with utils package

	// Handle graceful shutdown via SIGINT / SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-stop:
			// Gracefully shutdown the consumer
			utils.LogInfo("Received shutdown signal. Stopping consumer gracefully...")
			CloseConsumer(consumer) // Close the consumer gracefully
			return
		default:
			// Consume message with a small timeout (to avoid blocking forever)
			msg, err := consumer.ReadMessage(100 * time.Millisecond) // Timeout after 100ms
			if err != nil {
				// Log error but continue consuming
				kafkaError, ok := err.(kafka.Error)
				if ok && kafkaError.Code() != kafka.ErrTimedOut {
					// Only log actual errors
					utils.LogError("Error while consuming message", err)  // Logging with utils package
				}
				continue // Continue consuming even if there's an error
			}

			// Log the consumed message
			utils.LogInfo(fmt.Sprintf("Consumed message from topic %s: key = %-10s value = %s", *msg.TopicPartition.Topic, string(msg.Key), string(msg.Value)))  // Logging with utils package
			// You can add more message processing logic here
		}
	}
}

// CloseConsumer gracefully shuts down the Kafka consumer
func CloseConsumer(consumer *kafka.Consumer) {
	if err := consumer.Close(); err != nil {
		utils.LogError("Error while closing Kafka consumer", err)  // Logging with utils package
		return
	}
	utils.LogInfo("Kafka consumer closed successfully")
}
