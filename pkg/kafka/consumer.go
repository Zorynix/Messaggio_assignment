package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, groupID, topic string) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &KafkaConsumer{
		reader: r,
	}
}

func (kc *KafkaConsumer) Consume(ctx context.Context, handler func(ctx context.Context, message string)) {
	run := true
	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			m, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				if err == context.DeadlineExceeded || err == context.Canceled {
					logrus.Debugf("Consumer timeout: %v", err)
					continue
				}
				logrus.Errorf("Consumer error: %v", err)
				continue
			}
			handler(ctx, string(m.Value))
		}
	}

	if err := kc.reader.Close(); err != nil {
		logrus.Errorf("Failed to close consumer: %v", err)
	}
}

func (kp *KafkaConsumer) Close() {
	if err := kp.reader.Close(); err != nil {
		logrus.Errorf("Failed to close producer: %v", err)
	}
}
