package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	messagesMu sync.Mutex
	messages   []string
	brokerAddr = getEnv("KAFKA_BROKER", "kafkanew:9092")
	topic      = getEnv("KAFKA_TOPIC", "TEST_TOPIC")
)

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func main() {
	// Start a background goroutine to consume from Kafka
	go consumeMessages()

	// Setup HTTP routes
	http.HandleFunc("/messages", handleGetMessages) // GET /messages
	http.HandleFunc("/read", handleForceRead)       // POST /read

	addr := ":8081"
	log.Printf("Consumer HTTP server listening on %s, broker=%s, topic=%s", addr, brokerAddr, topic)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// consumeMessages runs forever, reading from Kafka, storing in memory
func consumeMessages() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddr},
		Topic:    topic,
		GroupID:  "consumerGroup", // or read from ENV
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer r.Close()

	log.Printf("Starting consumer for topic=%s, broker=%s", topic, brokerAddr)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		msg := string(m.Value)
		log.Printf("Received message: %s", msg)

		messagesMu.Lock()
		messages = append(messages, msg)
		messagesMu.Unlock()
	}
}

// handleGetMessages returns JSON array of messages
func handleGetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	messagesMu.Lock()
	defer messagesMu.Unlock()

	// Return the messages as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// handleForceRead triggers a manual read iteration (for demonstration)
func handleForceRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	// In this example, 'consumeMessages()' is already in a loop,
	// so there's nothing special to do. We'll just respond.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Consumer is reading messages in the background.\n"))
}

// func main() {
// 	// Kafka configuration
// 	brokers := []string{"kafkanew:9092"} // k8-staging
// 	// brokers := []string{"kafka:9092"} // multi container when using docker compose --build
// 	// brokers := []string{"localhost:29092"}
// 	topic := "TEST_TOPIC"
// 	groupID := "consumerGroup"

// 	// Create a new context
// 	ctx := context.Background()

// 	// Initialize a new reader with the brokers, topic, and groupID
// 	// The reader reads messages from the Kafka topic.
// 	reader := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:  brokers,
// 		Topic:    topic,
// 		GroupID:  groupID,
// 		MinBytes: 10e3, // 10KB
// 		MaxBytes: 10e6, // 10MB
// 	})

// 	defer reader.Close()

// 	fmt.Println("Start consuming ...")

// 	// Read messages
// 	for {
// 		msg, err := reader.ReadMessage(ctx)
// 		if err != nil {
// 			log.Fatalf("failed to read message: %s", err)
// 		}
// 		fmt.Printf("Message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))

// 		// Simulate processing time
// 		// time.Sleep(1 * time.Second)
// 	}
// }
