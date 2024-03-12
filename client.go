package messaging

type Serializer interface {
	Serialize(any) ([]byte, error)
	Deserialize([]byte, any) error
}

type ConfigProvider interface {
	GetBrokersAddrs() []string
}

type Client interface {
	Connect(ConfigProvider) error
	Send([]byte) error
	Receive() []byte
}
