package messaging

import "context"

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

func (m *MessageReceiver) ReceiveMessage(ctx context.Context) (Message, error) {
	message, err := m.client.Receive(ctx)
	if ctx.Err() != nil {
		return nil, err
	}
	if err != nil {
		return nil, NewErrReceive(err)
	}

	return message, nil
}

func (m *MessageReceiver) Close() error {
	return m.client.Close()
}
