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
	t.Run("connects to client", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigA{
			brokers: []string{"192.168.0.1"},
		}

		_, err := NewMessageSender(client, config, nil)

		AssertEqual(t, err, nil)
		AssertEqual(t, client.isConnected, true)
	})

	t.Run("returns ErrConnect on connection error", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigB{
			brokers:   []string{"192.168.0.1"},
			partition: 2,
		}

		_, err := NewMessageSender(client, config, nil)

		AssertErrorType[*ErrConnect](t, err)
	})
	t.Run("encodes message in JSON and sends it via SenderFunc", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigA{}
		serializer := &JSONSerializer{}

		sender, err := NewMessageSender(client, config, serializer)
		AssertEqual(t, err, nil)

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

		sender, err := NewMessageSender(client, config, serializer)
		AssertEqual(t, err, nil)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		msgpack.Unmarshal(client.data, &got)

		AssertEqual(t, got, want)
	})
}
