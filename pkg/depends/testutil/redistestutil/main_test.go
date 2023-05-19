package redistestutil_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/testutil/redistestutil"
)

func TestInit(t *testing.T) {
	NewWithT(t).Expect(Endpoint.Endpoint.String()).To(Equal("tcp://127.0.0.1:6379"))
	NewWithT(t).Expect(Endpoint.Prefix).To(Equal("prod:test:"))

	NewWithT(t).Expect(Redis.Host).To(Equal("127.0.0.1"))
	NewWithT(t).Expect(Redis.Port).To(Equal(6379))
}

func DISABLE_TestLivenessCheck(t *testing.T) {
	kvs := Redis.LivenessCheck()
	NewWithT(t).Expect(kvs[Redis.Host]).To(Equal("ok"))

	kvs = Endpoint.LivenessCheck()
	NewWithT(t).Expect(kvs[Endpoint.Endpoint.Host()]).To(Equal("ok"))
}
