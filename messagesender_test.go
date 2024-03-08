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

type StubConfigA struct {
	paramInt int
	param2   string
	param3   string
}

type StubClientA struct {
	config StubConfigA
	data   []byte
}

func (s *StubClientA) Connect(config ConnectConfig) error {
	s.config = StubConfigA{
		paramInt: config.paramInt,
		param2:   config.param2,
		param3:   config.param3,
	}

	return nil
}

func (s *StubClientA) Send(data []byte) error {
	s.data = data
	return nil
}

type StubConfigB struct {
	param1    string
	paramInt  int
	paramBool bool
}

type StubClienntB struct {
	config StubConfigB
	data   []byte
}

func (s *StubClienntB) Connect(config ConnectConfig) error {
	s.config = StubConfigB{
		param1:    config.param1,
		paramInt:  config.paramInt,
		paramBool: config.paramBool,
	}

	return nil
}

func (s *StubClienntB) Send(data []byte) error {
	s.data = data
	return nil
}

func TestMessageSender(t *testing.T) {

	t.Run("configures client and connects to broker", func(t *testing.T) {
		clientA := &StubClientA{}
		config := ConnectConfig{
			param1:    "Param 1",
			param2:    "Param 2",
			param3:    "Param 3",
			paramInt:  42,
			paramBool: true,
		}

		sender := NewMessageSender(clientA, config, json.Marshal)
		// bypass unused compiler error
		sender.marshalFunc("")

		want := StubConfigA{
			paramInt: config.paramInt,
			param2:   config.param2,
			param3:   config.param3,
		}
		AssertEqual(t, clientA.config, want)
	})

	t.Run("configures client B and connects to broker", func(t *testing.T) {
		clientB := &StubClienntB{}
		config := ConnectConfig{
			param1:    "Param 1",
			param2:    "Param 2",
			param3:    "Param 3",
			paramInt:  42,
			paramBool: true,
		}

		sender := NewMessageSender(clientB, config, json.Marshal)
		// bypass unused compiler error
		sender.marshalFunc("")

		want := StubConfigB{
			param1:    config.param1,
			paramInt:  config.paramInt,
			paramBool: config.paramBool,
		}
		AssertEqual(t, clientB.config, want)
	})

	t.Run("encodes message in JSON and sends it via SenderFunc", func(t *testing.T) {
		client := &StubClientA{}
		sender := NewMessageSender(client, ConnectConfig{}, json.Marshal)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		json.Unmarshal(client.data, &got)

		AssertEqual(t, got, want)
	})

	t.Run("encodes message in MessagePack and sends it via SenderFunc", func(t *testing.T) {
		client := &StubClientA{}
		sender := NewMessageSender(client, ConnectConfig{}, msgpack.Marshal)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		var got BaseMessage
		msgpack.Unmarshal(client.data, &got)

		AssertEqual(t, got, want)
	})
}
