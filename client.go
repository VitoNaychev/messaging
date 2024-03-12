package messaging

type ConfigProvider interface {
	GetBrokersAddrs() []string
}

type Client interface {
	Connect(ConfigProvider) error
	Send([]byte) error
	Receive() []byte
}
