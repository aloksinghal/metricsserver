package publish

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer kafka.Writer
}

func NewKafkaPublisher(brokerUrl string, topic string) Publisher {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerUrl},
		Topic:   topic,
		Balancer: &kafka.LeastBytes{},
	})
	return KafkaPublisher{
		writer:*w,
	}
}

func (k KafkaPublisher) Publish(messages [][]byte, key string) error {
	fmt.Println("received call to publish")
	var err error
	for _, message := range messages {
		err = k.writer.WriteMessages(context.Background(), kafka.Message{
			Key:  []byte(key),
			Value: message,
		})
	}
	return err
}