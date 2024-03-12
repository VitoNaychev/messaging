package messaging

type MarshalFunc func(any) ([]byte, error)

type MessageSender struct {
	client      Client
	marshalFunc MarshalFunc
}

func NewMessageSender(client Client, configProvider ConfigProvider, marshalFunc MarshalFunc) *MessageSender {
	messageSender := MessageSender{
		client:      client,
		marshalFunc: marshalFunc,
	}

	messageSender.client.Connect(configProvider)
	return &messageSender
}

func (m *MessageSender) SendMessage(message Message) error {
	data, _ := m.marshalFunc(message)
	return m.client.Send(data)
}
