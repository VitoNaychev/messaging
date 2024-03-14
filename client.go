package messaging

type ConfigProvider interface {
	GetBrokersAddrs() []string
}

type SenderConfigProvider interface {
	GetBrokersAddrs() []string
}

type ReceiverConfigProvider interface {
	GetBrokersAddrs() []string
	GetTopic() string
}

type Client interface {
	Connect(ConfigProvider) error
	Send(Message) error
	Receive() (Message, error)
}
