package config

import (
	"github.com/spf13/viper"
	"log"
)

// LoadConfig loads configuration from a file (e.g., config.yaml or config.json)
func LoadConfig() {
	viper.SetConfigName("config")  // Name of the configuration file (without extension)
	viper.AddConfigPath(".")       // Path to look for the config file
	viper.AutomaticEnv()           // Automatically read environment variables if set

	// Set defaults for Kafka-related configuration
	viper.SetDefault("kafka.bootstrapServers", "pkc-12576z.us-west2.gcp.confluent.cloud:9092")
	viper.SetDefault("kafka.apiKey", "S55GEJ2ZEI3UPWCR")
	viper.SetDefault("kafka.apiSecret", "6HMjpvY/o+i/+1s50polJAB2UDnksVj/Xk6gHlcl+C6iFYIV4HHq/PIQnmjME2I2")
	viper.SetDefault("kafka.topic", "navi-dev-topic")
	viper.SetDefault("server.port", "8080")

	// Read in the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Debug log to check the value of kafka.bootstrapServers
	log.Printf("Kafka Bootstrap Servers: %v", viper.GetString("kafka.bootstrapServers"))
}

// GetKafkaBootstrapServers returns the Kafka brokers from the config
func GetKafkaBootstrapServers() string {
	return viper.GetString("kafka.bootstrapServers")
}

// GetKafkaAPIKey returns the Kafka API Key from the config
func GetKafkaAPIKey() string {
	return viper.GetString("kafka.apiKey")
}

// GetKafkaAPISecret returns the Kafka API Secret from the config
func GetKafkaAPISecret() string {
	return viper.GetString("kafka.apiSecret")
}

// GetKafkaTopic returns the Kafka topic from the config
func GetKafkaTopic() string {
	return viper.GetString("kafka.topic")
}

// GetServerPort returns the server port from the config
func GetServerPort() string {
	return viper.GetString("server.port")
}
