package messaging

type MessageSender struct {
	client SenderClient
}

func NewMessageSender(client SenderClient, configProvider SenderConfigProvider) (*MessageSender, error) {
	messageSender := MessageSender{
		client: client,
	}

	err := messageSender.client.Connect(configProvider)
	if err != nil {
		return &messageSender, NewErrConnect(err)
	}

	return &messageSender, nil
}

func (m *MessageSender) SendMessage(message Message) error {
	err := m.client.Send(message)
	if err != nil {
		return NewErrSend(err)
	}

	return nil
}
