package messaging

type MarshalFunc func(any) ([]byte, error)

type Client interface {
	Connect() error
	Send([]byte) error
}

type MessageSender struct {
	client      Client
	marshalFunc MarshalFunc
}

func NewMessageSender(client Client, marshalFunc MarshalFunc) *MessageSender {
	messageSender := MessageSender{
		client:      client,
		marshalFunc: marshalFunc,
	}

	messageSender.client.Connect()
	return &messageSender
}

func (m *MessageSender) SendMessage(message Message) error {
	data, _ := m.marshalFunc(message)
	return m.client.Send(data)
}
