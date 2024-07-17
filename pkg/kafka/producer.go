package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(brokers string, topic string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: p,
		topic:    topic,
	}, nil
}

func (kp *KafkaProducer) Produce(message string) error {
	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	err := kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan
	msg, ok := e.(*kafka.Message)
	if !ok {
		return fmt.Errorf("unexpected event type %T", e)
	}
	if msg.TopicPartition.Error != nil {
		return msg.TopicPartition.Error
	}
	log.Infof("Message delivered to %v", msg.TopicPartition)
	return nil
}

func (kp *KafkaProducer) Close() {
	kp.producer.Close()
}
