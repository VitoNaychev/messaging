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
		serializer := &JSONSerializer{}

		receiver, err := NewMessageReceiver(client, config, serializer)
		AssertEqual(t, err, nil)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		client.data, _ = json.Marshal(want)

		got, err := receiver.ReceiveMessage()
		AssertEqual(t, err, nil)

		AssertEqual(t, got, (Message)(want))
	})

	t.Run("receives MessagePack encoded message via client", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigA{
			brokers: []string{"192.168.0.1"},
		}
		serializer := &MsgpackSerializer{}

		receiver, err := NewMessageReceiver(client, config, serializer)
		AssertEqual(t, err, nil)

		want := NewBaseMessage(testMessageID, testTopic, testPayload)
		client.data, _ = msgpack.Marshal(want)

		got, err := receiver.ReceiveMessage()
		AssertEqual(t, err, nil)

		AssertEqual(t, got, (Message)(want))
	})

	t.Run("returns ErrReceive on error during message receival", func(t *testing.T) {
		client := &StubClientA{}
		config := &StubConfigA{
			brokers: []string{"192.168.0.1"},
		}
		serializer := &MsgpackSerializer{}

		receiver, err := NewMessageReceiver(client, config, serializer)
		AssertEqual(t, err, nil)

		message := NewBaseMessage(testMessageID, testTopic, testPayload)
		client.data, _ = msgpack.Marshal(message)

		client.err = errors.New("dummy error")
		_, err = receiver.ReceiveMessage()
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
