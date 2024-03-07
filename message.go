package messaging

import (
	"time"

	"github.com/google/uuid"
)

type Message interface {
	GetTimestamp() time.Time
	GetUUID() uuid.UUID
	GetPayload() interface{}
}

type BaseMessage struct {
	Timestamp time.Time
	UUID      uuid.UUID
	Payload   interface{}
}

func (m *BaseMessage) GetTimestamp() time.Time {
	return m.Timestamp
}

func (m *BaseMessage) GetUUID() uuid.UUID {
	return m.UUID
}

func (m *BaseMessage) GetPayload() interface{} {
	return m.Payload
}
