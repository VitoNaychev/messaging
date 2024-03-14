package messaging

type ReceiverConfigProvider interface {
	GetBrokersAddrs() []string
	GetTopic() string
}

type ReceiverClient interface {
	Connect(ReceiverConfigProvider) error
	Receive() (Message, error)
}
