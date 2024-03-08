package messaging

import (
	"encoding/json"
	"testing"
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
	sender := NewMessageSender(messageSender.SenderFunc)

	t.Run("sends message to channel", func(t *testing.T) {
		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		json.Unmarshal(messageSender.data, &got)

		AssertEqual(t, want, got)
	})
}
