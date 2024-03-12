package messaging

type StubConfigA struct {
	brokers []string
}

func (s *StubConfigA) GetBrokersAddrs() []string {
	return s.brokers
}

type StubClientA struct {
	brokers []string
	data    []byte
}

func (s *StubClientA) Connect(config ConfigProvider) error {
	s.brokers = config.GetBrokersAddrs()

	return nil
}

func (s *StubClientA) Send(data []byte) error {
	s.data = data
	return nil
}

type StubConfigB struct {
	brokers []string
}

func (s *StubConfigB) GetBrokersAddrs() []string {
	return s.brokers
}

type StubClientB struct {
	brokers []string
	data    []byte
}

func (s *StubClientB) Connect(config ConfigProvider) error {
	s.brokers = config.GetBrokersAddrs()

	return nil
}

func (s *StubClientB) Send(data []byte) error {
	s.data = data
	return nil
}
