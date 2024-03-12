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

func TestMessageSender(t *testing.T) {
	t.Run("encodes message in JSON and sends it via SenderFunc", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigA{}
		serializer := &JSONSerializer{}
		sender := NewMessageSender(client, config, serializer)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		json.Unmarshal(client.data, &got)

		AssertEqual(t, got, want)
	})

	t.Run("encodes message in MessagePack and sends it via SenderFunc", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigA{}
		serializer := &MsgpackSerializer{}
		sender := NewMessageSender(client, config, serializer)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		msgpack.Unmarshal(client.data, &got)

		AssertEqual(t, got, want)
	})
}
