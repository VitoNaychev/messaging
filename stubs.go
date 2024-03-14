package messaging

import (
	"errors"
)

var (
	ErrConfigMismatch = errors.New("Client doesn't support this ConfigProvider")
)

type StubConfigA struct {
	brokers []string
	topic   string

	connType string
}

func (s *StubConfigA) GetBrokersAddrs() []string {
	return s.brokers
}

func (s *StubConfigA) GetTopic() string {
	return s.topic
}

func (s *StubConfigA) GetConnectionType() string {
	return s.connType
}

type StubClientA struct {
	isConnected bool
	err         error

	brokers  []string
	topic    string
	connType string

	message Message
}

func (s *StubClientA) Connect(config ConfigProvider) error {
	s.isConnected = true

	configA, ok := config.(*StubConfigA)
	if !ok {
		return ErrConfigMismatch
	}

	s.brokers = configA.GetBrokersAddrs()

	if config, ok := config.(ReceiverConfigProvider); ok {
		s.topic = config.GetTopic()
	}

	s.connType = configA.GetConnectionType()

	return nil
}

func (s *StubClientA) Send(message Message) error {
	s.message = message
	return s.err
}

func (s *StubClientA) Receive() (Message, error) {
	return s.message, s.err
}

type StubConfigB struct {
	brokers []string
	topic   string

	partition int
}

func (s *StubConfigB) GetBrokersAddrs() []string {
	return s.brokers
}

func (s *StubConfigB) GetTopic() string {
	return s.topic
}

func (s *StubConfigB) GetPartition() int {
	return s.partition
}

type StubClientB struct {
	brokers   []string
	partition int

	message Message
}

func (s *StubClientB) Connect(config ConfigProvider) error {
	configB, ok := config.(*StubConfigB)
	if !ok {
		return ErrConfigMismatch
	}

	s.brokers = configB.GetBrokersAddrs()
	s.partition = configB.GetPartition()

	return nil
}

func (s *StubClientB) Send(message Message) error {
	s.message = message
	return nil
}

func (s *StubClientB) Receive() (Message, error) {
	return s.message, nil
}
