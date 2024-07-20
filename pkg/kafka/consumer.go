package kafka

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    string
}

func NewKafkaConsumer(brokers, groupID, topic string) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: c,
		topic:    topic,
	}, nil
}

func (kc *KafkaConsumer) Consume(ctx context.Context, handler func(ctx context.Context, message string)) {
	run := true
	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			msg, err := kc.consumer.ReadMessage(time.Hour)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
					logrus.Debugf("Consumer timeout: %v", err)
					continue
				}
				logrus.Errorf("Consumer error: %v", err)
				continue
			}

			if msg != nil && msg.Value != nil {
				message := string(msg.Value)
				handler(ctx, message)
			} else {
				logrus.Warn("Received an empty message or nil value")
			}
		}
	}

	kc.consumer.Close()
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}
