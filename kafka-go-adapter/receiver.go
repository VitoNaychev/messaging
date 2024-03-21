package kafkagoadapter

import (
	"context"
	"encoding/json"
	"errors"
	"messaging"

	"github.com/segmentio/kafka-go"
)

type KafkaReceiverConfigProvider struct {
	ConsumerGroup string
	messaging.BaseReceiverConfigProvider
}

func (k *KafkaReceiverConfigProvider) GetConsumerGroup() string {
	return k.ConsumerGroup
}

type Receiver struct {
	reader *kafka.Reader
}

func (r *Receiver) Connect(config messaging.ReceiverConfigProvider) error {
	kafkaConfig, ok := config.(*KafkaReceiverConfigProvider)
	if !ok {
		return errors.New("client doesn't support this ConfigProvider")
	}

	r.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kafkaConfig.GetBrokersAddrs(),
		GroupID:  kafkaConfig.GetConsumerGroup(),
		Topic:    kafkaConfig.GetTopic(),
		MaxBytes: 10e6, // 10MB
	})

	return nil
}

func (r *Receiver) Receive(ctx context.Context) (messaging.Message, error) {
	kafkaMessage, err := r.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	message := messaging.BaseMessage{}
	json.Unmarshal(kafkaMessage.Value, &message)
	return message, nil
}

func (s *Receiver) Close() error {
	return s.reader.Close()
}
