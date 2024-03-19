package messaging

import (
	"context"
)

type MessageHandler func(Message) error

type MessageRouter struct {
	client      ReceiverClient
	subscribers map[int]MessageHandler
	errChan     chan error
}

func NewMessageRouter(client ReceiverClient, config ReceiverConfigProvider) (*MessageRouter, error) {
	router := &MessageRouter{
		client:      client,
		subscribers: map[int]MessageHandler{},
		errChan:     make(chan error, 1),
	}

	err := router.client.Connect(config)
	if err != nil {
		return nil, NewErrConnect(err)
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
		message, err := m.client.Receive(ctx)
		if ctx.Err() != nil {
			return err
		}

		err = m.subscribers[message.GetMessageID()](message)
		if err != nil {
			m.errChan <- err
		}
	}
}

func (m *MessageRouter) Errors() chan error {
	return m.errChan
}
