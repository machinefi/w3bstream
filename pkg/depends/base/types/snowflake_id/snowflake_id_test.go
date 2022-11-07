package snowflake_id_test

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/base/types/snowflake_id"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
)

func TestWorkerIDFromIP(t *testing.T) {
	t.Log(WorkerIDFromIP(net.ParseIP("255.255.255.255")))
	t.Log(WorkerIDFromIP(net.ParseIP("127.0.0.1")))
	t.Log(WorkerIDFromLocalIP())
}

func NewSnowflakeTestSuite(t *testing.T, n int) *SnowflakeTestSuite {
	return &SnowflakeTestSuite{t, n, mapx.New[uint64, bool]()}
}

type SnowflakeTestSuite struct {
	*testing.T
	N int
	*mapx.Map[uint64, bool]
}

func (s *SnowflakeTestSuite) ExpectN(n int) {
	NewWithT(s.T).Expect(s.Len()).To(Equal(n))
}

func (s *SnowflakeTestSuite) Run(sf *Snowflake) {
	for i := 1; i <= s.N; i++ {
		id, err := sf.ID()
		if err != nil {
			s.T.Log(err)
		}
		NewWithT(s.T).Expect(err).To(BeNil())
		s.Store(id, true)
	}
}

func BenchmarkSnowflake_ID(b *testing.B) {
	start, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	cases := [][3]uint{
		{11, 11, 1},
		{10, 12, 1},
		{12, 10, 1},
		{12, 12, 1},
		{16, 10, 1},
		{16, 8, 1},
		{16, 8, 5},
		{16, 8, 10},
	}

	for _, vs := range cases {
		f := NewSnowflakeFactory(vs[0], vs[1], vs[2], start.Local())
		s, err := f.NewSnowflake(1)
		if err != nil {
			b.Fatal(err)
		}
		name := fmt.Sprintf(
			"END_%s__MAX_WORKER_%d_MAX_SEQ_%d_PER_%dms",
			f.MaxUnitTime(), f.MaxWorkerID(), f.MaxSequence(), vs[2],
		)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err = s.ID()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func TestSnowflake(t *testing.T) {
	suite := NewSnowflakeTestSuite(t, 1000000)
	g, err := NewSnowflake(1)
	NewWithT(t).Expect(err).To(BeNil())

	suite.Run(g)
	suite.ExpectN(suite.N)
}

func TestSnowflake_Concurrent(t *testing.T) {
	suite := NewSnowflakeTestSuite(t, 100)
	g, err := NewSnowflake(1)
	NewWithT(t).Expect(err).To(BeNil())

	con := 1000
	wg := &sync.WaitGroup{}

	for i := 0; i < con; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			suite.Run(g)
		}()
	}

	wg.Wait()
	suite.ExpectN(suite.N * con)
}

func TestSnowflakeFactory_BuildIDx(t *testing.T) {
	f, err := NewSnowflakeFactory(10, 12, 1, time.Now().In(time.UTC)).NewSnowflake(294)
	NewWithT(t).Expect(err).To(BeNil())

	for i := 0; i < 10000; i++ {
		_, err := f.ID()
		NewWithT(t).Expect(err).To(BeNil())
		time.Sleep(100 * time.Microsecond)
	}
}
