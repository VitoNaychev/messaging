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

type UnmarshalFunc func([]byte, any) error

type MessageReceiver struct {
	unmarshal UnmarshalFunc
	client    Client
}

func NewMessageReceiver(client Client, configProvider ConfigProvider, unmarshal UnmarshalFunc) (*MessageReceiver, error) {
	receiver := MessageReceiver{
		unmarshal: unmarshal,
		client:    client,
	}
	err := receiver.client.Connect(configProvider)
	if err != nil {
		return nil, NewErrConnect(err)
	}
	return &receiver, nil
}

func (m *MessageReceiver) ReceiveMessage() Message {
	var message BaseMessage
	m.unmarshal(m.client.Receive(), &message)

	return message
}
