package messaging

type SenderConfigProvider interface {
	GetBrokersAddrs() []string
}

type BaseSenderConfigProvider struct {
	Brokers []string
}

func (s *BaseSenderConfigProvider) GetBrokersAddrs() []string {
	return s.Brokers
}

type SenderClient interface {
	Connect(SenderConfigProvider) error
	Send(Message) error
	Close() error
}
