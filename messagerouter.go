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
	client ReceiverClient

	subscribers map[int]MessageHandler
	errors      bool
	errChan     chan error
}

func NewMessageRouter(client ReceiverClient, config *RouterConfigProvider) (*MessageRouter, error) {
	router := &MessageRouter{
		client: client,

		subscribers: map[int]MessageHandler{},
		errors:      config.Errors,
		errChan:     make(chan error, 1),
	}

	err := router.client.Connect(config.ReceiverConfigProvider)
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

func (m *MessageRouter) Errors() chan error {
	return m.errChan
}
