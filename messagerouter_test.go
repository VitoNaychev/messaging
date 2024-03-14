package messaging

import (
	"errors"
	"reflect"
	"testing"
)

var dummyMessageID = 10

func dummyMessageHandler(message Message) error {
	return nil
}

func TestMessageRouter(t *testing.T) {
	t.Run("connects to client", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
			topic:   "test-topic",
		}

		_, err := NewMessageRouter(client, config)

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
		_, err := NewMessageRouter(client, config)

		AssertErrorType[*ErrConnect](t, err)
	})

	t.Run("adds message handler to list of subscribers", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
			topic:   "test-topic",
		}

		router, err := NewMessageRouter(client, config)
		AssertEqual(t, err, nil)

		router.Subscribe(dummyMessageID, dummyMessageHandler)

		got := router.subscribers[dummyMessageID]
		want := dummyMessageHandler
		AssertEqualFunc(t, got, want)
	})
}

func AssertEqualFunc[T any](t testing.TB, got T, want T) {
	t.Helper()

	if reflect.TypeOf(got).Kind() != reflect.Func {
		t.Errorf("expected func arguments got %v", reflect.TypeOf(got).Kind())
	}

	gotType := reflect.TypeOf(got)
	wantType := reflect.TypeOf(want)

	if gotType != wantType {
		t.Errorf("got function type %v want %v", gotType, wantType)
	}

	gotAddress := reflect.ValueOf(got)
	wantAddress := reflect.ValueOf(want)

	if gotAddress != wantAddress {
		t.Errorf("got function address %v want %v", gotAddress, wantAddress)
	}
}
