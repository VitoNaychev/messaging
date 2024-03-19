package messaging

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

var dummyError = errors.New("dummy error")

var dummyMessageID = 10

type StubMessageHandler struct {
	message Message
	err     error
}

func (s *StubMessageHandler) Handle(message Message) error {
	s.message = message
	return s.err
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

		messageHandler := &StubMessageHandler{}
		err = router.Subscribe(dummyMessageID, messageHandler.Handle)
		AssertEqual(t, err, nil)

		got := router.subscribers[dummyMessageID]
		want := messageHandler.Handle
		AssertEqualFunc(t, got, want)
	})

	t.Run("returns ErrDuplicateHandler if a handlers is already registered", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
			topic:   "test-topic",
		}

		router, err := NewMessageRouter(client, config)
		AssertEqual(t, err, nil)

		messageHandler := &StubMessageHandler{}
		err = router.Subscribe(dummyMessageID, messageHandler.Handle)
		AssertEqual(t, err, nil)

		err = router.Subscribe(dummyMessageID, messageHandler.Handle)
		AssertEqual(t, err, ErrDuplicateHandler)
	})

	t.Run("listens until context is done", func(t *testing.T) {
		client := &StubReceiverClient{
			timeout: time.Millisecond,
		}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
			topic:   "test-topic",
		}

		router, err := NewMessageRouter(client, config)
		AssertEqual(t, err, nil)

		messageHandler := &StubMessageHandler{}
		err = router.Subscribe(dummyMessageID, messageHandler.Handle)
		AssertEqual(t, err, nil)

		want := NewBaseMessage(dummyMessageID, "test-topic", "Hello, World!")
		client.message = want

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2)
		defer cancel()

		err = router.Listen(ctx)
		AssertEqual(t, err, ctx.Err())
	})

	t.Run("routes message to message handler", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
			topic:   "test-topic",
		}

		router, err := NewMessageRouter(client, config)
		AssertEqual(t, err, nil)

		messageHandler := &StubMessageHandler{}
		router.Subscribe(dummyMessageID, messageHandler.Handle)

		want := NewBaseMessage(dummyMessageID, "test-topic", "Hello, World!")
		client.message = want

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		router.Listen(ctx)

		got := messageHandler.message
		AssertEqual(t, got, (Message)(want))
	})

	t.Run("forwards handler errors via Errors() chan", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &StubConfig{
			brokers: []string{"192.168.0.1"},
			topic:   "test-topic",
		}

		router, err := NewMessageRouter(client, config)
		AssertEqual(t, err, nil)

		messageHandler := &StubMessageHandler{err: dummyError}
		router.Subscribe(dummyMessageID, messageHandler.Handle)

		want := NewBaseMessage(dummyMessageID, "test-topic", "Hello, World!")
		client.message = want

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		router.Listen(ctx)

		err = <-router.Errors()
		AssertEqual(t, err, messageHandler.err)
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

	gotAddress := reflect.ValueOf(got).Pointer()
	wantAddress := reflect.ValueOf(want).Pointer()

	if gotAddress != wantAddress {
		t.Errorf("got function address %v want %v", gotAddress, wantAddress)
	}
}
