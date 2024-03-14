package messaging

type SenderConfigProvider interface {
	GetBrokersAddrs() []string
}

type SenderClient interface {
	Connect(SenderConfigProvider) error
	Send(Message) error
}
