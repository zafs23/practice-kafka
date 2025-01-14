package main

import (
	"fmt"
	"net"

	"github.com/segmentio/kafka-go"
)

func createTopic(brokers []string, topic string) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return err
	}

	defer conn.Close()

	// zookeeper --> 3 brokers, and there is a leader which is the controller
	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	// now we need a connection to the controller and give an address

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, fmt.Sprint(controller.Port)))
	if err != nil {
		return err
	}

	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	return controllerConn.CreateTopics(topicConfig)
}

/*

 1 2 3

*/
