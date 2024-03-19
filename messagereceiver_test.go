package messaging

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestMessageReceiver(t *testing.T) {
	t.Run("connects to client", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
			topic:   "test-topic",
		}

		_, err := NewMessageReceiver(client, config)

		AssertEqual(t, err, nil)
		AssertEqual(t, client.isConnected, true)
	})

	t.Run("returns ErrConnect on connection error", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
			topic:   "test-topic",
		}

		client.err = errors.New("dummy error")
		_, err := NewMessageReceiver(client, config)

		AssertErrorType[*ErrConnect](t, err)
	})

	t.Run("receives message", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
			topic:   "test-topic",
		}

		receiver, err := NewMessageReceiver(client, config)
		AssertEqual(t, err, nil)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		client.message = want

		got, err := receiver.ReceiveMessage(context.Background())
		AssertEqual(t, err, nil)

		AssertEqual(t, got, (Message)(want))
	})

	t.Run("cancels receive message via context", func(t *testing.T) {
		client := &StubReceiverClient{
			timeout: time.Millisecond * 2,
		}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
			topic:   "test-topic",
		}

		receiver, err := NewMessageReceiver(client, config)
		AssertEqual(t, err, nil)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		client.message = want

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_, err = receiver.ReceiveMessage(ctx)
		AssertEqual(t, err, ctx.Err())

	})

	t.Run("returns ErrReceive on error during message receival", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
		}

		receiver, err := NewMessageReceiver(client, config)
		AssertEqual(t, err, nil)

		message := NewBaseMessage(testMessageID, testTopic, testPayload)
		client.message = message

		client.err = errors.New("dummy error")
		_, err = receiver.ReceiveMessage(context.Background())
		AssertErrorType[*ErrReceive](t, err)

	})

}

func AssertErrorType[T error](t testing.TB, got error) {
	var want T

	if !errors.As(got, &want) {
		t.Errorf("got error with type %v want %v", reflect.TypeOf(got), reflect.TypeOf(want))
	}
}

func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
