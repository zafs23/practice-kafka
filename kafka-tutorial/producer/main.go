package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	brokers := []string{"kafka:9092"} // broker list
	// brokers := []string{"localhost:29092"}
	// create a topic
	topic := "TEST_TOPIC"

	err := createTopic(brokers, topic)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Topic created")

	/**
	after creating a topic we need to write and need to know which brokers and topics we
	are writing to
	*/

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	ctx := context.Background()
	// write message will write a bach of messages
	err = writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World 2!"),
		})

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Successfully published message to topic")
}
