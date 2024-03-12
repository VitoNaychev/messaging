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

func (m *MessageReceiver) ReceiveMessage() (Message, error) {
	data, err := m.client.Receive()
	if err != nil {
		return nil, NewErrReceive(err)
	}

	var message BaseMessage
	err = m.serializer.Deserialize(data, &message)
	if err != nil {
		return nil, NewErrReceive(err)
	}

	return message, nil
}
