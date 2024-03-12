package messaging

type MessageSender struct {
	client     Client
	serializer Serializer
}

func NewMessageSender(client Client, configProvider ConfigProvider, serializer Serializer) *MessageSender {
	messageSender := MessageSender{
		client:     client,
		serializer: serializer,
	}

	messageSender.client.Connect(configProvider)
	return &messageSender
}

func (m *MessageSender) SendMessage(message Message) error {
	data, _ := m.serializer.Serialize(message)
	return m.client.Send(data)
}
