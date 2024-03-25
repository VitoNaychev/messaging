package messaging

import (
	"context"
	"errors"
	"time"
)

var (
	ErrConfigMismatch = errors.New("client doesn't support this ConfigProvider")
)

type StubConfig struct {
	brokers []string
	topic   string

	connType string
}

func (s *StubConfig) GetBrokersAddrs() []string {
	return s.brokers
}

func (s *StubConfig) GetTopic() string {
	return s.topic
}

func (s *StubConfig) GetConnectionType() string {
	return s.connType
}

type StubSenderClient struct {
	isConnected bool
	isClosed    bool
	err         error

	brokers  []string
	connType string

	message Message
}

func (s *StubSenderClient) Connect(config SenderConfigProvider) error {
	s.isConnected = true

	stubConfig, ok := config.(*StubConfig)
	if !ok {
		return ErrConfigMismatch
	}

	s.brokers = stubConfig.GetBrokersAddrs()
	s.connType = stubConfig.GetConnectionType()

	return s.err
}

func (s *StubSenderClient) Send(message Message) error {
	s.message = message
	return s.err
}

func (s *StubSenderClient) Close() error {
	s.isClosed = true
	return nil
}

type StubReceiverClient struct {
	isConnected bool
	isClosed    bool
	timeout     time.Duration
	err         error

	brokers  []string
	topic    string
	connType string

	message Message
}

func (s *StubReceiverClient) Connect(config ReceiverConfigProvider) error {
	s.isConnected = true

	stubConfig, ok := config.(*StubConfig)
	if !ok {
		return ErrConfigMismatch
	}

	s.brokers = stubConfig.GetBrokersAddrs()
	s.topic = stubConfig.GetTopic()

	s.connType = stubConfig.GetConnectionType()

	return s.err
}

func (s *StubReceiverClient) Receive(ctx context.Context) (Message, error) {
	sleepChan := make(chan interface{})

	go func() {
		if s.timeout != 0 {
			time.Sleep(s.timeout)
		}
		sleepChan <- nil
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-sleepChan:
		return s.message, s.err
	}
}

func (s *StubReceiverClient) Close() error {
	s.isClosed = true
	return nil
}
