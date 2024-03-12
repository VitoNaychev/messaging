package messaging

type MessageReceiver struct {
	serializer Serializer
	client     Client
}

func NewMessageReceiver(client Client, configProvider ConfigProvider, serializer Serializer) (*MessageReceiver, error) {
	receiver := MessageReceiver{
		serializer: serializer,
		client:     client,
	}

	err := receiver.client.Connect(configProvider)
	if err != nil {
		return nil, NewErrConnect(err)
	}

	return &receiver, nil
}

func (m *MessageReceiver) ReceiveMessage() Message {
	var message BaseMessage
	data, _ := m.client.Receive()
	m.serializer.Deserialize(data, &message)

	return message
}
