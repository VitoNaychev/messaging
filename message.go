package messaging

import (
	"time"

	"github.com/google/uuid"
)

type Message interface {
	GetMessageID() int
	GetTopic() string
	GetTimestamp() time.Time
	GetUUID() string
	GetPayload() interface{}
}

type BaseMessage struct {
	UUID      string
	Timestamp time.Time
	Topic     string
	MessageID int
	Payload   interface{}
}

func NewBaseMessage(messageID int, topic string, payload interface{}) BaseMessage {
	return BaseMessage{
		UUID:      uuid.NewString(),
		Timestamp: time.Now(),
		Topic:     topic,
		MessageID: messageID,
		Payload:   payload,
	}
}

func (m BaseMessage) GetMessageID() int {
	return m.MessageID
}

func (m BaseMessage) GetTopic() string {
	return m.Topic
}

func (m BaseMessage) GetTimestamp() time.Time {
	return m.Timestamp
}

func (m BaseMessage) GetUUID() string {
	return m.UUID
}

func (m BaseMessage) GetPayload() interface{} {
	return m.Payload
}
