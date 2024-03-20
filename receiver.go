package messaging

import "context"

type ReceiverConfigProvider interface {
	GetBrokersAddrs() []string
	GetTopic() string
}

type ReceiverClient interface {
	Connect(ReceiverConfigProvider) error
	Receive(context.Context) (Message, error)
	Close() error
}
