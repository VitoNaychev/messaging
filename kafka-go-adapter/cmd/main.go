package main

import (
	"context"
	"fmt"
	"log"
	"messaging"
	kafkagoadapter "messaging/kafka-go-adapter"
	"time"
)

var messasgeID = 5
var topic = "my-topic"
var payload = "Hello, World"

func main() {
	senderConfig := &messaging.BaseSenderConfigProvider{
		Brokers: []string{"localhost:9092"},
	}
	sender, err := messaging.NewMessageSender(&kafkagoadapter.Sender{}, senderConfig)
	if err != nil {
		log.Fatal("NewMessageSender error: ", err)
	}

	sender.SendMessage(messaging.NewBaseMessage(messasgeID, topic, payload))

	receiverConfig := &kafkagoadapter.KafkaReceiverConfigProvider{
		ConsumerGroup: "my-consumer-group",
		BaseReceiverConfigProvider: messaging.BaseReceiverConfigProvider{
			Brokers: []string{"localhost:9092"},
			Topic:   topic,
		},
	}
	receiver, err := messaging.NewMessageReceiver(&kafkagoadapter.Receiver{}, receiverConfig)
	if err != nil {
		log.Fatal("NewMessageReceiver error: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	message, err := receiver.ReceiveMessage(ctx)
	if err != nil {
		log.Fatal("ReceiveMessage error: ", err)
	}

	fmt.Println(message)
}
