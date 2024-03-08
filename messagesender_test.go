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
	config ConnectConfig
	data   []byte
}

func (s *StubClient) Connect(config ConnectConfig) error {
	s.config = config

	return nil
}

func (s *StubClient) Send(data []byte) error {
	s.data = data
	return nil
}

func TestMessageSender(t *testing.T) {
	client := &StubClient{}

	t.Run("configures client and connects to broker", func(t *testing.T) {
		config := ConnectConfig{
			param1: "Param 1",
			param2: "Param 2",
			param3: "Param 3",
		}

		sender := NewMessageSender(client, config, json.Marshal)
		sender.marshalFunc("")

		AssertEqual(t, client.config, config)
	})

	t.Run("encodes message in JSON and sends it via SenderFunc", func(t *testing.T) {
		sender := NewMessageSender(client, ConnectConfig{}, json.Marshal)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		json.Unmarshal(client.data, &got)

		AssertEqual(t, got, want)
	})

	t.Run("encodes message in MessagePack and sends it via SenderFunc", func(t *testing.T) {
		sender := NewMessageSender(client, ConnectConfig{}, msgpack.Marshal)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		msgpack.Unmarshal(client.data, &got)

		AssertEqual(t, got, want)
	})
}
