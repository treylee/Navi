package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
    "navi-ears/config"
    "navi-ears/consumer"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
}

var clients = make(map[*websocket.Conn]bool)

// Handle WebSocket requests
func handleConnections(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()
    clients[conn] = true

    for {
        _, _, err := conn.ReadMessage() // We could receive messages, but for now just keep the connection open
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
        log.Println("Consumed message:", string(msg.Value))

        // Send consumed Kafka message to React clients via WebSocket
        broadcastMessage(string(msg.Value))
    }
}

func main() {
    // Start WebSocket server
    http.HandleFunc("/ws", handleConnections)
    go func() {
        log.Println("WebSocket server started on :8080")
        if err := http.ListenAndServe(":8080", nil); err != nil {
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

    consumeMessages(consumer)
}
