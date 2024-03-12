package messaging

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func TestMessageReceiver(t *testing.T) {
	t.Run("connects to client", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigA{
			brokers: []string{"192.168.0.1"},
		}

		_, err := NewMessageReceiver(client, config, nil)

		AssertEqual(t, err, nil)
		AssertEqual(t, client.isConnected, true)
	})

	t.Run("returns ErrConnect on connection error", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigB{
			brokers:   []string{"192.168.0.1"},
			partition: 2,
		}

		_, err := NewMessageReceiver(client, config, nil)

		AssertErrorType[*ErrConnect](t, err)
	})

	t.Run("receives JSON encoded message via client", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigA{
			brokers: []string{"192.168.0.1"},
		}

		receiver, err := NewMessageReceiver(client, config, json.Unmarshal)
		AssertEqual(t, err, nil)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		client.data, _ = json.Marshal(want)

		got := receiver.ReceiveMessage()

		AssertEqual(t, got, (Message)(want))
	})

	t.Run("receives MessagePack encoded message via client", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigA{
			brokers: []string{"192.168.0.1"},
		}

		receiver, err := NewMessageReceiver(client, config, msgpack.Unmarshal)
		AssertEqual(t, err, nil)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		client.data, _ = msgpack.Marshal(want)

		got := receiver.ReceiveMessage()

		AssertEqual(t, got, (Message)(want))
	})

	// t.Run("receives message from channel", func(t *testing.T) {
	// 	want := NewBaseMessage(testMessageID, testTopic, testPayload)

	// 	got := receiver.ReceiveMessage()

	// 	if !reflect.DeepEqual(want, got) {
	// 		t.Errorf("got %v want %v", got, want)
	// 	}
	// })

	// t.Run("receives multiple messages from channel", func(t *testing.T) {
	// 	wantFirst := NewBaseMessage(testMessageID, testTopic, "First")
	// 	wantSecond := NewBaseMessage(testMessageID, testTopic, "Second")
	// 	wantThird := NewBaseMessage(testMessageID, testTopic, "Third")

	// 	go func() {
	// 		msgChan <- wantFirst
	// 		msgChan <- wantSecond
	// 		msgChan <- wantThird
	// 	}()

	// 	got := receiver.ReceiveMessage()
	// 	AssertEqual(t, got, (Message)(wantFirst))

	// 	got = receiver.ReceiveMessage()
	// 	AssertEqual(t, got, (Message)(wantSecond))

	// 	got = receiver.ReceiveMessage()
	// 	AssertEqual(t, got, (Message)(wantThird))
	// })

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
