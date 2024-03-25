package messaging

import (
	"errors"
	"testing"
)

var (
	testMessageID = 1
	testTopic     = "test-topic"
	testPayload   = "Hello, World!"
)

func TestMessageSender(t *testing.T) {
	t.Run("connects to client", func(t *testing.T) {
		client := &StubSenderClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
		}

		_, err := NewMessageSender(client, config)

		AssertEqual(t, err, nil)
		AssertEqual(t, client.isConnected, true)
	})

	t.Run("returns ErrConnect on connection error", func(t *testing.T) {
		client := &StubSenderClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
		}

		client.err = errors.New("dummy error")
		_, err := NewMessageSender(client, config)

		AssertErrorType[*ErrConnect](t, err)
	})
	t.Run("sends message", func(t *testing.T) {
		client := &StubSenderClient{}
		config := &StubConfig{}

		sender, err := NewMessageSender(client, config)
		AssertEqual(t, err, nil)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		got := client.message

		AssertEqual(t, got, (Message)(want))
	})

	t.Run("returns ErrSend on error during message sending", func(t *testing.T) {
		client := &StubSenderClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
		}

		sender, err := NewMessageSender(client, config)
		AssertEqual(t, err, nil)

		client.err = errors.New("dummy error")
		message := NewBaseMessage(testMessageID, testTopic, testPayload)
		err = sender.SendMessage(message)

		AssertErrorType[*ErrSend](t, err)
	})

	t.Run("calls client.Close in Close method", func(t *testing.T) {
		client := &StubSenderClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
		}

		sender, err := NewMessageSender(client, config)
		AssertEqual(t, err, nil)

		sender.Close()
		AssertEqual(t, client.isClosed, true)
	})
}
