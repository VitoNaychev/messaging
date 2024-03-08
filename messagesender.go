package messaging

type MessageSender struct {
	msgChan chan Message
}

func NewMessageSender(msgChan chan Message) *MessageSender {
	return &MessageSender{
		msgChan: msgChan,
	}
}

func (m *MessageSender) SendMessage(message Message) {
	m.msgChan <- message
}
