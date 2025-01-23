package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

// MessagePayload is the JSON structure for incoming POST requests
type MessagePayload struct {
	Message string `json:"message"`
}

var (
	// We'll read the broker address from ENV or default
	brokerAddr = getEnv("KAFKA_BROKER", "kafkanew:9092")
	topic      = getEnv("KAFKA_TOPIC", "TEST_TOPIC")
	// Create a Kafka writer at startup
	kafkaWriter *kafka.Writer
)

// getEnv is a helper to read environment variables with a fallback
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func main() {
	// Initialize a Kafka writer
	kafkaWriter = &kafka.Writer{
		Addr:         kafka.TCP(brokerAddr),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}
	defer kafkaWriter.Close()

	// Setup HTTP routes
	http.HandleFunc("/send", handleSendMessage)

	// Start an HTTP server on port 8080
	addr := ":8080"
	log.Printf("Producer HTTP server listening on %s, broker=%s, topic=%s", addr, brokerAddr, topic)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// handleSendMessage handles POST /send with JSON {"message":"..."}
func handleSendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload MessagePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	msg := strings.TrimSpace(payload.Message)
	if msg == "" {
		http.Error(w, "Message cannot be empty", http.StatusBadRequest)
		return
	}

	// Produce the message to Kafka
	err := kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(fmt.Sprintf("key-%v", msg)),
			Value: []byte(msg),
		},
	)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		http.Error(w, "Failed to write to Kafka", http.StatusInternalServerError)
		return
	}

	log.Printf("Produced message to topic=%s: %s", topic, msg)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Message sent: %s\n", msg)))
}
