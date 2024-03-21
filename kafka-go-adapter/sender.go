package kafkagoadapter

import (
	"context"
	"encoding/json"
	"messaging"

	"github.com/segmentio/kafka-go"
)

type Sender struct {
	writer *kafka.Writer
}

func (s *Sender) Connect(config messaging.SenderConfigProvider) error {
	s.writer = &kafka.Writer{
		Addr:                   kafka.TCP(config.GetBrokersAddrs()...),
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.LeastBytes{},
	}

	return nil
}

func (s *Sender) Send(message messaging.Message) error {
	msgJSON, _ := json.Marshal(message)

	kafkaMessage := kafka.Message{
		Topic: message.GetTopic(),
		Value: msgJSON,
	}
	return s.writer.WriteMessages(context.Background(), kafkaMessage)
}

func (s *Sender) Close() error {
	return s.writer.Close()
}
