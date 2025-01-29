package config

import (
	"github.com/spf13/viper"
	"log"
)

// LoadConfig loads configuration from config files (e.g., config.json, config.yaml)
func LoadConfig() {
	viper.SetConfigName("config")  // Name of the configuration file (without extension)
	viper.AddConfigPath(".")       // Path to look for the config file

	viper.AutomaticEnv() // Automatically read environment variables if set

	// Set defaults if the config is missing
	viper.SetDefault("kafka.bootstrapServers", "pkc-12576z.us-west2.gcp.confluent.cloud:9092")
	viper.SetDefault("kafka.apiKey", "S55GEJ2ZEI3UPWCR")
	viper.SetDefault("kafka.apiSecret", "6HMjpvY/o+i/+1s50polJAB2UDnksVj/Xk6gHlcl+C6iFYIV4HHq/PIQnmjME2I2") // Replace with your actual API secret
	viper.SetDefault("kafka.topic", "navi-dev-topic")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

// GetKafkaBootstrapServers returns the Kafka bootstrap server address
func GetKafkaBootstrapServers() string {
	return viper.GetString("kafka.bootstrapServers")
}

// GetKafkaAPIKey returns the Kafka API Key
func GetKafkaAPIKey() string {
	return viper.GetString("kafka.apiKey")
}

// GetKafkaAPISecret returns the Kafka API Secret
func GetKafkaAPISecret() string {
	return viper.GetString("kafka.apiSecret")
}

// GetKafkaTopic returns the Kafka topic to consume from
func GetKafkaTopic() string {
	return viper.GetString("kafka.topic")
}
