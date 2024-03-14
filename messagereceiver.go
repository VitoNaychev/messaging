package messaging

type MessageReceiver struct {
	client Client
}

func NewMessageReceiver(client Client, configProvider ReceiverConfigProvider) (*MessageReceiver, error) {
	receiver := MessageReceiver{
		client: client,
	}

	err := receiver.client.Connect(configProvider)
	if err != nil {
		return nil, NewErrConnect(err)
	}

	return &receiver, nil
}

func (m *MessageReceiver) ReceiveMessage() (Message, error) {
	message, err := m.client.Receive()
	if err != nil {
		return nil, NewErrReceive(err)
	}

	return message, nil
}
