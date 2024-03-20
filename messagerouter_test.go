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
		config := &RouterConfigProvider{
			Errors: false,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
		}

		_, err := NewMessageRouter(client, config)

		AssertEqual(t, err, nil)
		AssertEqual(t, client.isConnected, true)
	})

	t.Run("returns ErrConnect on connection error", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &RouterConfigProvider{
			Errors: false,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
		}

		client.err = errors.New("dummy error")
		_, err := NewMessageRouter(client, config)

		AssertErrorType[*ErrConnect](t, err)
	})

	t.Run("adds message handler to list of subscribers", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &RouterConfigProvider{
			Errors: false,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
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
		config := &RouterConfigProvider{
			Errors: false,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
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
		config := &RouterConfigProvider{
			Errors: false,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
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

	t.Run("returns ErrReceive on error during message receival", func(t *testing.T) {
		client := &StubReceiverClient{
			timeout: time.Millisecond,
		}
		config := &RouterConfigProvider{
			Errors: false,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
		}

		router, err := NewMessageRouter(client, config)
		AssertEqual(t, err, nil)

		messageHandler := &StubMessageHandler{}
		err = router.Subscribe(dummyMessageID, messageHandler.Handle)
		AssertEqual(t, err, nil)

		want := NewBaseMessage(dummyMessageID, "test-topic", "Hello, World!")
		client.message = want

		client.err = errors.New("dummy error")
		err = router.Listen(context.Background())
		AssertErrorType[*ErrReceive](t, err)
	})

	t.Run("routes message to message handler", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &RouterConfigProvider{
			Errors: false,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
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

	t.Run("forwards handler errors on config.Errors = true", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &RouterConfigProvider{
			Errors: true,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
		}

		router, err := NewMessageRouter(client, config)
		AssertEqual(t, err, nil)

		messageHandler := &StubMessageHandler{err: dummyError}
		router.Subscribe(dummyMessageID, messageHandler.Handle)

		want := NewBaseMessage(dummyMessageID, "test-topic", "Hello, World!")
		client.message = want

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		go listenForErrors(ctx, router, &err)
		router.Listen(ctx)

		AssertEqual(t, err, messageHandler.err)
	})

	t.Run("doesn't forward handler errors on config.Errors = false", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &RouterConfigProvider{
			Errors: false,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
		}

		router, err := NewMessageRouter(client, config)
		AssertEqual(t, err, nil)

		messageHandler := &StubMessageHandler{err: dummyError}
		router.Subscribe(dummyMessageID, messageHandler.Handle)

		want := NewBaseMessage(dummyMessageID, "test-topic", "Hello, World!")
		client.message = want

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		go listenForErrors(ctx, router, &err)
		router.Listen(ctx)

		AssertEqual(t, err, nil)
	})

	t.Run("sends ErrUnknownMessage on unknown message and config.Errors = true", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &RouterConfigProvider{
			Errors: true,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
		}

		router, err := NewMessageRouter(client, config)
		AssertEqual(t, err, nil)

		messageHandler := &StubMessageHandler{err: dummyError}
		router.Subscribe(dummyMessageID, messageHandler.Handle)

		unknownMessageID := 42
		want := NewBaseMessage(unknownMessageID, "test-topic", "Hello, World!")
		client.message = want

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		go listenForErrors(ctx, router, &err)
		router.Listen(ctx)

		AssertEqual(t, err, ErrUnknownMessage)
	})

	t.Run("skips message handler on unknown message and config.Errors = false", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &RouterConfigProvider{
			Errors: false,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
		}

		router, err := NewMessageRouter(client, config)
		AssertEqual(t, err, nil)

		messageHandler := &StubMessageHandler{err: dummyError}
		router.Subscribe(dummyMessageID, messageHandler.Handle)

		unknownMessageID := 42
		want := NewBaseMessage(unknownMessageID, "test-topic", "Hello, World!")
		client.message = want

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		go listenForErrors(ctx, router, &err)
		router.Listen(ctx)

		AssertEqual(t, err, nil)
	})

	t.Run("calls client.Close in Close method", func(t *testing.T) {
		client := &StubReceiverClient{}
		config := &RouterConfigProvider{
			Errors: false,
			ReceiverConfigProvider: &StubConfig{
				brokers: []string{"192.168.0.1"},
				topic:   "test-topic",
			},
		}

		router, err := NewMessageRouter(client, config)
		AssertEqual(t, err, nil)

		router.Close()
		AssertEqual(t, client.isClosed, true)
	})
}

func listenForErrors(ctx context.Context, router *MessageRouter, err *error) {
	for {
		select {
		case *err = <-router.Errors():
		case <-ctx.Done():
			return
		}
	}
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
