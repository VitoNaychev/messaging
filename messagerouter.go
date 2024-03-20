package messaging

import (
	"context"
)

type RouterConfigProvider struct {
	Errors bool
	ReceiverConfigProvider
}

type MessageHandler func(Message) error

type MessageRouter struct {
	receiver *MessageReceiver

	subscribers map[int]MessageHandler
	errors      bool
	errChan     chan error
}

func NewMessageRouter(client ReceiverClient, config *RouterConfigProvider) (*MessageRouter, error) {
	recevier, err := NewMessageReceiver(client, config.ReceiverConfigProvider)
	if err != nil {
		return nil, err
	}

	router := &MessageRouter{
		receiver: recevier,

		subscribers: map[int]MessageHandler{},
		errors:      config.Errors,
		errChan:     make(chan error, 1),
	}

	return router, err
}

func (m *MessageRouter) Subscribe(messageID int, handler MessageHandler) error {
	if _, ok := m.subscribers[messageID]; ok {
		return ErrDuplicateHandler
	}

	m.subscribers[messageID] = handler
	return nil
}

func (m *MessageRouter) Listen(ctx context.Context) error {
	for {
		message, err := m.receiver.ReceiveMessage(ctx)
		if err != nil {
			return err
		}

		handler, ok := m.subscribers[message.GetMessageID()]
		if !ok {
			if m.errors {
				m.errChan <- ErrUnknownMessage
			}
			continue
		}

		err = handler(message)
		if m.errors && err != nil {
			m.errChan <- err
		}
	}
}

func (m *MessageRouter) Close() error {
	return m.receiver.Close()
}

func (m *MessageRouter) Errors() chan error {
	return m.errChan
}
