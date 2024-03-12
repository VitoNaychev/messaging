package messaging

import (
	"testing"
)

func TestClient(t *testing.T) {
	t.Run("configures client and connects to broker", func(t *testing.T) {
		clientA := &StubClientA{}
		configA := &StubConfigA{
			brokers:  []string{"192.168.0.1", "192.168.0.255"},
			connType: "tcp",
		}

		err := clientA.Connect(configA)
		AssertEqual(t, err, nil)

		AssertEqual(t, clientA.brokers, configA.brokers)
		AssertEqual(t, clientA.connType, configA.connType)
	})

	t.Run("configures client B and connects to broker", func(t *testing.T) {
		clientB := &StubClientB{}
		configB := &StubConfigB{
			brokers:   []string{"192.168.0.1"},
			partition: 2,
		}

		err := clientB.Connect(configB)
		AssertEqual(t, err, nil)

		AssertEqual(t, clientB.brokers, configB.brokers)
		AssertEqual(t, clientB.partition, configB.partition)
	})

	t.Run("returns ErrConfigMismatch on wrong config", func(t *testing.T) {
		clientB := &StubClientB{}
		configA := &StubConfigA{
			brokers:  []string{"192.168.0.1"},
			connType: "tcp",
		}

		err := clientB.Connect(configA)
		// bypass unused compiler error

		AssertEqual(t, err, ErrConfigMismatch)
	})
}
