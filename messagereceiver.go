package messaging

type MessageReceiver struct {
	msgChan chan Message
}

func NewMessageReceiver(msgChan chan Message) *MessageReceiver {
	return &MessageReceiver{
		msgChan: msgChan,
	}
}

func (m *MessageReceiver) ReceiveMessage() Message {
	return <-m.msgChan
}
