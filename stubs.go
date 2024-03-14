package messaging

import (
	"errors"
)

var (
	ErrConfigMismatch = errors.New("Client doesn't support this ConfigProvider")
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

func (s *StubSenderClient) Receive() (Message, error) {
	return s.message, s.err
}

type StubReceiverClient struct {
	isConnected bool
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

func (s *StubReceiverClient) Send(message Message) error {
	s.message = message
	return s.err
}

func (s *StubReceiverClient) Receive() (Message, error) {
	return s.message, s.err
}
