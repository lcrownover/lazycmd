package main

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lcrownover/lazycmd/internal/lazycmd"
)

func main() {

	myHostname := lazycmd.GetHostname()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          myHostname,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.Subscribe(myHostname, nil)

	// can use this to break out of the forever
	run := true

	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on: %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if err.(kafka.Error).Code() != kafka.ErrTimedOut {
			// The client will automatically try to recover from all errors.
			// kafka.ErrTimedOut is not considered an error because it is
			// raised by ReadMessage on timeout.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}
