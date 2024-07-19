package kafka

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
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

func (kc *KafkaConsumer) Consume(ctx context.Context) {
	run := true

	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			msg, err := kc.consumer.ReadMessage(time.Second)
			if err == nil {
				log.Infof("Message on %s: %s", msg.TopicPartition, string(msg.Value))
			} else if err.(kafka.Error).IsTimeout() {
				log.Debugf("Consumer timeout: %v", err)
			} else {
				log.Errorf("Consumer error: %v (%v)", err, msg)
			}
		}
	}

	kc.consumer.Close()
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}
