package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/alex-pvl/go-tapmenu/internal/app/config"
	"github.com/alex-pvl/go-tapmenu/internal/app/store"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(config *config.Configuration) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.KafkaAddress),
		RequiredAcks: -1,
		MaxAttempts:  5,
		BatchSize:    200,
		WriteTimeout: 30 * time.Second,
		Balancer:     &kafka.RoundRobin{},
	}

	return &Producer{writer: writer}
}

func (p *Producer) SendMessage(message store.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messageJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	kafkaMsg := kafka.Message{
		Topic: message.RestaurantName,
		Value: messageJson,
	}

	return p.writer.WriteMessages(ctx, kafkaMsg)
}
