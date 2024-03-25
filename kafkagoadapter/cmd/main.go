package main

import (
	"context"
	"fmt"
	"log"
	"messaging"
	"messaging/kafkagoadapter"
	"time"
)

var greetingMessageID = 5
var farewellMessageID = 10
var topic = "my-topic"

func GreetingMessageHandler(message messaging.Message) error {
	fmt.Println("received a greeting: ", message.GetPayload())
	return nil
}

func FarewellMessageHandler(message messaging.Message) error {
	fmt.Println("received a farewell: ", message.GetPayload())
	return nil
}

func main() {
	senderConfig := &messaging.BaseSenderConfigProvider{
		Brokers: []string{"localhost:9092"},
	}
	sender, err := messaging.NewMessageSender(&kafkagoadapter.SenderClient{}, senderConfig)
	if err != nil {
		log.Fatal("NewMessageSender error: ", err)
	}

	sender.SendMessage(messaging.NewBaseMessage(greetingMessageID, topic, "Hello, World!"))
	sender.SendMessage(messaging.NewBaseMessage(farewellMessageID, topic, "Goodbye, World :("))

	routerConfig := &messaging.RouterConfigProvider{
		Errors: true,
		ReceiverConfigProvider: &kafkagoadapter.KafkaReceiverConfigProvider{
			ConsumerGroup: "my-consumer-group",
			BaseReceiverConfigProvider: messaging.BaseReceiverConfigProvider{
				Brokers: []string{"localhost:9092"},
				Topic:   topic,
			},
		},
	}
	router, err := messaging.NewMessageRouter(&kafkagoadapter.ReceiverClient{}, routerConfig)
	if err != nil {
		log.Fatal("NewMessageReceiver error: ", err)
	}

	router.Subscribe(greetingMessageID, GreetingMessageHandler)
	router.Subscribe(farewellMessageID, FarewellMessageHandler)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	router.Listen(ctx)
}
