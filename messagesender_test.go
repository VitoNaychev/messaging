package messaging

import (
	"reflect"
	"testing"
)

var (
	testMessageID = 1
	testTopic     = "test-topic"
	testPayload   = "Hello, World!"
)

func TestMessageSender(t *testing.T) {
	msgChan := make(chan Message)
	sender := NewMessageSender(msgChan)

	t.Run("sends message to channel", func(t *testing.T) {
		var got Message
		go func() { got = <-msgChan }()

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		sender.SendMessage(want)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
