package messaging

import "errors"

var (
	ErrConfigMismatch = errors.New("Client doesn't support this ConfigProvider")
)

type StubConfigA struct {
	brokers  []string
	connType string
}

func (s *StubConfigA) GetBrokersAddrs() []string {
	return s.brokers
}

func (s *StubConfigA) GetConnectionType() string {
	return s.connType
}

type StubClientA struct {
	isConnected bool

	brokers  []string
	connType string

	data []byte
}

func (s *StubClientA) Connect(config ConfigProvider) error {
	s.isConnected = true

	configA, ok := config.(*StubConfigA)
	if !ok {
		return ErrConfigMismatch
	}

	s.brokers = configA.GetBrokersAddrs()
	s.connType = configA.GetConnectionType()

	return nil
}

func (s *StubClientA) Send(data []byte) error {
	s.data = data
	return nil
}

type StubConfigB struct {
	brokers   []string
	partition int
}

func (s *StubConfigB) GetBrokersAddrs() []string {
	return s.brokers
}

func (s *StubConfigB) GetPartition() int {
	return s.partition
}

type StubClientB struct {
	brokers   []string
	partition int
	data      []byte
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

func (s *StubClientB) Send(data []byte) error {
	s.data = data
	return nil
}
