package messaging

type MarshalFunc func(any) ([]byte, error)
type SenderFunc func([]byte) error

type MessageSender struct {
	senderFunc  SenderFunc
	marshalFunc MarshalFunc
}

func NewMessageSender(senderFunc SenderFunc, marshalFunc MarshalFunc) *MessageSender {
	return &MessageSender{
		senderFunc:  senderFunc,
		marshalFunc: marshalFunc,
	}
}

func (m *MessageSender) SendMessage(message Message) error {
	data, _ := m.marshalFunc(message)
	return m.senderFunc(data)
}
