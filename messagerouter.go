package messaging

type MessageHandler func(Message) error

type MessageRouter struct {
	client      ReceiverClient
	subscribers map[int]MessageHandler
}

func NewMessageRouter(client ReceiverClient, config ReceiverConfigProvider) (*MessageRouter, error) {
	router := &MessageRouter{
		client:      client,
		subscribers: map[int]MessageHandler{},
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

func (m *MessageRouter) Listen() {

}
