package messaging

import (
	"encoding/json"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

var (
	testMessageID = 1
	testTopic     = "test-topic"
	testPayload   = "Hello, World!"
)

type StubSenderFunc struct {
	data []byte
}

func (s *StubSenderFunc) SenderFunc(data []byte) error {
	s.data = data
	return nil
}

func TestMessageSender(t *testing.T) {
	messageSender := StubSenderFunc{}

	t.Run("encodes message in JSON and sends it via SenderFunc", func(t *testing.T) {
		sender := NewMessageSender(messageSender.SenderFunc, json.Marshal)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		json.Unmarshal(messageSender.data, &got)

		AssertEqual(t, got, want)
	})

	t.Run("encodes message in MessagePack and sends it via SenderFunc", func(t *testing.T) {
		sender := NewMessageSender(messageSender.SenderFunc, msgpack.Marshal)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		msgpack.Unmarshal(messageSender.data, &got)

		AssertEqual(t, got, want)
	})
}
