package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
    "navi-ears/config"
    "navi-ears/consumer"
    "os"
    "os/signal"
    "syscall"
    "time"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
}

var clients = make(map[*websocket.Conn]bool) // Track WebSocket clients

// Handle WebSocket requests
func handleConnections(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()

    clients[conn] = true
    log.Printf("New WebSocket connection: %s", conn.RemoteAddr())

    for {
        _, _, err := conn.ReadMessage() // Keep the connection open (could be used to receive messages later)
        if err != nil {
            log.Println("Read error:", err)
            delete(clients, conn)
            break
        }
    }
}

// Broadcast messages to all connected WebSocket clients
func broadcastMessage(message string) {
    for client := range clients {
        err := client.WriteMessage(websocket.TextMessage, []byte(message))
        if err != nil {
            log.Println("WebSocket send error:", err)
            client.Close()
            delete(clients, client)
        } else {
            log.Printf("Message sent to WebSocket client: %s", client.RemoteAddr())
        }
    }
}

// Kafka message consumer
func consumeMessages(consumer *kafka.Consumer) {
    for {
        msg, err := consumer.ReadMessage(-1)
        if err != nil {
            log.Println("Error consuming message:", err)
            continue
        }

        // Log the Kafka message value before sending to WebSocket
        log.Printf("Consumed message: %s", string(msg.Value))

        // Send consumed Kafka message to React clients via WebSocket
        broadcastMessage(string(msg.Value))
    }
}

func main() {
	config.LoadConfig()
    // Handle graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    // Start WebSocket server
    http.HandleFunc("/ws", handleConnections)
    go func() {
        log.Println("WebSocket server started on :8083")
        if err := http.ListenAndServe(":8083", nil); err != nil {
            log.Fatal("ListenAndServe error:", err)
        }
    }()

    // Start Kafka consumer
    consumer, err := consumer.CreateConsumer()
    if err != nil {
        log.Fatalf("Error creating Kafka consumer: %v", err)
    }
    defer consumer.Close()

    topic := config.GetKafkaTopic()
    err = consumer.Subscribe(topic, nil)
    if err != nil {
        log.Fatalf("Error subscribing to Kafka topic: %v", err)
    }

    go consumeMessages(consumer)

    log.Println("Consumer started. Press CTRL+C to stop.")

    // Wait for shutdown signal
    <-stop
    log.Println("Received shutdown signal. Shutting down gracefully...")

    // Optional: Timeout to wait for WebSocket connections to close
    shutdownTimeout := time.After(10 * time.Second)
    select {
    case <-shutdownTimeout:
        log.Println("Shutdown timeout reached. Forcefully shutting down...")
    }

    // Close WebSocket connections (you may want to close active connections more gracefully)
    for client := range clients {
        client.Close()
    }

    log.Println("Consumer and WebSocket server shut down successfully.")
}
