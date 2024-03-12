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

func TestClientInterface(t *testing.T) {
	t.Run("configures client and connects to broker", func(t *testing.T) {
		clientA := &StubClientA{}
		configA := &StubConfigA{
			brokers: []string{"192.168.0.1", "192.168.0.255"},
		}

		sender := NewMessageSender(clientA, configA, json.Marshal)
		// bypass unused compiler error
		sender.marshalFunc("")

		AssertEqual(t, clientA.brokers, configA.brokers)
	})

	t.Run("configures client B and connects to broker", func(t *testing.T) {
		clientB := &StubClientB{}
		configB := &StubConfigB{
			brokers: []string{"192.168.0.1"},
		}

		sender := NewMessageSender(clientB, configB, json.Marshal)
		// bypass unused compiler error
		sender.marshalFunc("")

		AssertEqual(t, clientB.brokers, configB.brokers)
	})
}

func TestMessageSender(t *testing.T) {
	t.Run("encodes message in JSON and sends it via SenderFunc", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigA{}
		sender := NewMessageSender(client, config, json.Marshal)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		json.Unmarshal(client.data, &got)

		AssertEqual(t, got, want)
	})

	t.Run("encodes message in MessagePack and sends it via SenderFunc", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigA{}
		sender := NewMessageSender(client, config, msgpack.Marshal)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		msgpack.Unmarshal(client.data, &got)

		AssertEqual(t, got, want)
	})
}
