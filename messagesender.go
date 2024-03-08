package messaging

import "encoding/json"

type SenderFunc func([]byte) error

type MessageSender struct {
	senderFunc SenderFunc
}

func NewMessageSender(senderFunc SenderFunc) *MessageSender {
	return &MessageSender{
		senderFunc: senderFunc,
	}
}

func (m *MessageSender) SendMessage(message Message) error {
	data, _ := json.Marshal(message)
	return m.senderFunc(data)
}
