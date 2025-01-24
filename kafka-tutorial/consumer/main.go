package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
)

// Prometheus metrics
var (
	messagesConsumed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_consumed_total",
			Help: "Total number of Kafka messages consumed",
		},
		[]string{"topic"},
	)
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

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(messagesConsumed)
}

func main() {
	// Start a background goroutine to consume from Kafka
	go consumeMessages()

	// Start Prometheus HTTP server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Starting Prometheus metrics server on :8082")
		log.Fatal(http.ListenAndServe(":8082", nil))
	}()

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

		// Increment Prometheus counter
		messagesConsumed.WithLabelValues(topic).Inc()

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
