package messaging

type MessageHandler interface {
	HandleMessage(message Message) error
}
