package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lcrownover/lazycmd/internal/lazycmd"
)

func main() {

	fTarget := flag.String("target", "", "target host")
	fCommand := flag.String("command", "", "command to run")
	flag.Parse()

	if *fTarget == "" || *fCommand == "" {
		fmt.Println("You must provide both a -target and -command")
		os.Exit(1)
	}

	target := lazycmd.CleanseTarget(*fTarget)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery Failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// produce messages asynchronously
	topic := target
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(*fCommand),
	}, nil)

	// wait for deliveries before shutting down
	p.Flush(15 * 1000)
}
