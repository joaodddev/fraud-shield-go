package kafka

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func NewProducer(broker, topic string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	slog.Info("kafka producer connected", "broker", broker, "topic", topic)

	return &Producer{
		producer: p,
		topic:    topic,
	}, nil
}

func (p *Producer) Publish(key string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: data,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	p.producer.Flush(3000)

	slog.Info("event published", "topic", p.topic, "key", key)
	return nil
}

func (p *Producer) Close() {
	p.producer.Close()
	slog.Info("kafka producer closed")
}
