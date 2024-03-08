package messaging

type ConnectConfig struct {
	param1 string
	param2 string
	param3 string
}

type MarshalFunc func(any) ([]byte, error)

type Client interface {
	Connect(ConnectConfig) error
	Send([]byte) error
}

type MessageSender struct {
	client      Client
	marshalFunc MarshalFunc
}

func NewMessageSender(client Client, config ConnectConfig, marshalFunc MarshalFunc) *MessageSender {
	messageSender := MessageSender{
		client:      client,
		marshalFunc: marshalFunc,
	}

	messageSender.client.Connect(config)
	return &messageSender
}

func (m *MessageSender) SendMessage(message Message) error {
	data, _ := m.marshalFunc(message)
	return m.client.Send(data)
}
