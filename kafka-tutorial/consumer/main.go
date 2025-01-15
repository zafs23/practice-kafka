package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Kafka configuration
	brokers := []string{"kafka:9092"} // multi container when using docker compose --build
	// brokers := []string{"localhost:29092"}
	topic := "TEST_TOPIC"
	groupID := "consumerGroup"

	// Create a new context
	ctx := context.Background()

	// Initialize a new reader with the brokers, topic, and groupID
	// The reader reads messages from the Kafka topic.
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	defer reader.Close()

	fmt.Println("Start consuming ...")

	// Read messages
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("failed to read message: %s", err)
		}
		fmt.Printf("Message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))

		// Simulate processing time
		// time.Sleep(1 * time.Second)
	}
}
