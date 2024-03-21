package kafkagoadapter

import (
	"context"
	"encoding/json"
	"messaging"

	"github.com/segmentio/kafka-go"
)

type SenderClient struct {
	writer *kafka.Writer
}

func (s *SenderClient) Connect(config messaging.SenderConfigProvider) error {
	s.writer = &kafka.Writer{
		Addr:                   kafka.TCP(config.GetBrokersAddrs()...),
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.LeastBytes{},
	}

	return nil
}

func (s *SenderClient) Send(message messaging.Message) error {
	msgJSON, _ := json.Marshal(message)

	kafkaMessage := kafka.Message{
		Topic: message.GetTopic(),
		Value: msgJSON,
	}
	return s.writer.WriteMessages(context.Background(), kafkaMessage)
}

func (s *SenderClient) Close() error {
	return s.writer.Close()
}
