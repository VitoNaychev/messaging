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

type StubClient struct {
	data []byte
}

func (s *StubClient) Connect() error {
	return nil
}

func (s *StubClient) Send(data []byte) error {
	s.data = data
	return nil
}

func TestMessageSender(t *testing.T) {
	client := &StubClient{}

	t.Run("encodes message in JSON and sends it via SenderFunc", func(t *testing.T) {
		sender := NewMessageSender(client, json.Marshal)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		json.Unmarshal(client.data, &got)

		AssertEqual(t, got, want)
	})

	t.Run("encodes message in MessagePack and sends it via SenderFunc", func(t *testing.T) {
		sender := NewMessageSender(client, msgpack.Marshal)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		msgpack.Unmarshal(client.data, &got)

		AssertEqual(t, got, want)
	})
}
