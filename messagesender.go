package messaging

type MessageSender struct {
	client     Client
	serializer Serializer
}

func NewMessageSender(client Client, configProvider ConfigProvider, serializer Serializer) (*MessageSender, error) {
	messageSender := MessageSender{
		client:     client,
		serializer: serializer,
	}

	err := messageSender.client.Connect(configProvider)
	if err != nil {
		return &messageSender, NewErrConnect(err)
	}

	return &messageSender, nil
}

func (m *MessageSender) SendMessage(message Message) error {
	data, _ := m.serializer.Serialize(message)
	return m.client.Send(data)
}
