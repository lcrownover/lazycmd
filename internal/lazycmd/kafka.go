package lazycmd

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	Topic    string
	Consumer *kafka.Consumer
}

func NewConsumer(topic string) *Consumer {
	rx, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          topic,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	rx.SubscribeTopics([]string{topic}, nil)
	c := &Consumer{
		Topic:    topic,
		Consumer: rx,
	}
	fmt.Printf("Subscribed to topic: %s\n", c.Topic)
	return c
}

func (c *Consumer) GetMessage() (string, bool) {
	msg, err := c.Consumer.ReadMessage(time.Second)
	if err == nil {
		// fmt.Printf("Message on: %s: %s\n", msg.TopicPartition, string(msg.Value))
		return string(msg.Value), true
	}
	switch err.(kafka.Error).Code() {
	case kafka.ErrUnknownTopicOrPart:
		// fmt.Printf("Subscribed to new topic: %s\n", c.Topic)
		return "", false
	case kafka.ErrTimedOut:
		return "", false
	default:
		// fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		return "", false
	}
}

func (c *Consumer) Close() {
	c.Consumer.Close()
}

type Producer struct {
	Topic    string
	Producer *kafka.Producer
}

func NewProducer(topic string) *Producer {
	tx, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})
	if err != nil {
		panic(err)
	}
	p := &Producer{
		Topic:    topic,
		Producer: tx,
	}
	p.WatchEvents()
	return p
}

func (p *Producer) SendMessage(topic string, message string) {
	p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
	p.Flush()
}

func (p *Producer) Events() chan kafka.Event {
	return p.Producer.Events()
}

func (p *Producer) WatchEvents() {
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Send failed for command: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Sent command to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
}

func (p *Producer) Flush() {
	p.Producer.Flush(15 * 1000)
}

func (p *Producer) Close() {
	p.Producer.Close()
}
