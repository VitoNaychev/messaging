package messaging

type MessageReceiver struct {
	client ReceiverClient
}

func NewMessageReceiver(client ReceiverClient, configProvider ReceiverConfigProvider) (*MessageReceiver, error) {
	messageReceiver := MessageReceiver{
		client: client,
	}

	err := messageReceiver.client.Connect(configProvider)
	if err != nil {
		return nil, NewErrConnect(err)
	}

	return &messageReceiver, nil
}

func (m *MessageReceiver) ReceiveMessage() (Message, error) {
	message, err := m.client.Receive()
	if err != nil {
		return nil, NewErrReceive(err)
	}

	return message, nil
}
