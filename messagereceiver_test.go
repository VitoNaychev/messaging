package messaging

import (
	"reflect"
	"testing"
)

func TestMessageReceiver(t *testing.T) {
	msgChan := make(chan Message)
	receiver := NewMessageReceiver(msgChan)

	t.Run("receives message from channel", func(t *testing.T) {
		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		go func() { msgChan <- want }()

		got := receiver.ReceiveMessage()

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("receives multiple messages from channel", func(t *testing.T) {
		wantFirst := NewBaseMessage(testMessageID, testTopic, "First")
		wantSecond := NewBaseMessage(testMessageID, testTopic, "Second")
		wantThird := NewBaseMessage(testMessageID, testTopic, "Third")

		go func() {
			msgChan <- wantFirst
			msgChan <- wantSecond
			msgChan <- wantThird
		}()

		got := receiver.ReceiveMessage()
		AssertEqual(t, got, (Message)(wantFirst))

		got = receiver.ReceiveMessage()
		AssertEqual(t, got, (Message)(wantSecond))

		got = receiver.ReceiveMessage()
		AssertEqual(t, got, (Message)(wantThird))
	})

}

func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
