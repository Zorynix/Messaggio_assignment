package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	return &KafkaProducer{
		writer: w,
	}
}

func (kp *KafkaProducer) Produce(ctx context.Context, message string) error {
	err := kp.writer.WriteMessages(ctx,
		kafka.Message{
			Value: []byte(message),
		},
	)
	if err != nil {
		logrus.Errorf("Failed to produce message: %v", err)
		return err
	}
	logrus.Infof("Message delivered to topic %v", kp.writer.Stats().Topic)
	return nil
}

func (kp *KafkaProducer) Close() {
	if err := kp.writer.Close(); err != nil {
		logrus.Errorf("Failed to close producer: %v", err)
	}
}
