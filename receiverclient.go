package messaging

import "context"

type ReceiverConfigProvider interface {
	GetBrokersAddrs() []string
	GetTopic() string
}

type BaseReceiverConfigProvider struct {
	Brokers []string
	Topic   string
}

func (k *BaseReceiverConfigProvider) GetBrokersAddrs() []string {
	return k.Brokers
}

func (k *BaseReceiverConfigProvider) GetTopic() string {
	return k.Topic
}

type ReceiverClient interface {
	Connect(ReceiverConfigProvider) error
	Receive(context.Context) (Message, error)
	Close() error
}
