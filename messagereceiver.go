package messaging

import (
	"fmt"
)

type ErrConnect struct {
	msg string
	err error
}

func NewErrConnect(err error) *ErrConnect {
	return &ErrConnect{
		msg: "got error during connection",
		err: err,
	}
}

func (e *ErrConnect) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e *ErrConnect) Unwrap() error {
	return e.err
}

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
	m.serializer.Deserialize(m.client.Receive(), &message)

	return message
}
