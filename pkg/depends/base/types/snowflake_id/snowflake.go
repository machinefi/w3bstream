package snowflake_id

import (
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrOverTimeLimit      = errors.New("over the time limit")
	ErrOverMaxWorkerID    = errors.New("over the worker id limit")
	ErrInvalidSystemClock = errors.New("invalid system clock")
)

func NewSnowflakeFactory(bitsWorkerID, bitsSequence, gap uint, base time.Time) *SnowflakeFactory {
	sf := &SnowflakeFactory{
		bitsWorkerID:  bitsWorkerID,
		bitsSequence:  bitsSequence,
		bitsTimestamp: 63 - bitsWorkerID - bitsSequence,
		unit:          time.Duration(gap) * time.Millisecond,
		base:          base,
	}
	sf.init()
	return sf
}

type SnowflakeFactory struct {
	bitsWorkerID  uint
	bitsSequence  uint
	bitsTimestamp uint
	unit          time.Duration
	base          time.Time

	maxWorkerID uint32    // maxWorkerID  = 1<<bitsWorkerID - 1
	maxSequence uint32    // maxSequence  = 1<<bitsSequence - 1
	maxUnitTime time.Time // maxUnitTime  =
}

func (f *SnowflakeFactory) init() {
	f.maxSequence = 1<<f.bitsSequence - 1
	f.maxWorkerID = 1<<f.bitsWorkerID - 1

	maxTS, baseTS := uint64(1<<f.bitsTimestamp-1), f.SnowFlakeTimestamp(f.base)
	seconds := int64(time.Duration(baseTS+maxTS) * f.unit / time.Second)
	f.maxUnitTime = time.Unix(seconds, 0).In(time.UTC)
}

func (f *SnowflakeFactory) MaxWorkerID() uint32 { return f.maxWorkerID }

func (f *SnowflakeFactory) MaxSequence() uint32 { return f.maxSequence }

func (f *SnowflakeFactory) MaxUnitTime() time.Time { return f.maxUnitTime }

func (f *SnowflakeFactory) MaskSequence(seq uint32) uint32 { return seq & f.maxSequence }

func (f *SnowflakeFactory) SnowFlakeTimestamp(t time.Time) uint64 {
	return uint64(t.In(time.UTC).UnixNano() / int64(f.unit))
}

func (f *SnowflakeFactory) Sleep(d time.Duration) time.Duration {
	return d*f.unit - time.Duration(time.Now().In(time.UTC).UnixNano())%f.unit*time.Nanosecond
}

func (f *SnowflakeFactory) Elapsed() uint64 {
	return f.SnowFlakeTimestamp(time.Now().In(time.UTC)) - f.SnowFlakeTimestamp(f.base.In(time.UTC))
}

func (f *SnowflakeFactory) BuildID(wid, seq uint32, elapsed uint64) (uint64, error) {
	if elapsed >= 1<<f.bitsTimestamp {
		return 0, ErrOverTimeLimit
	}
	v := elapsed<<(f.bitsSequence+f.bitsWorkerID) |
		uint64(seq)<<f.bitsWorkerID |
		uint64(wid)
	// log.Printf("build %d w: %d seq: %2d elapsed: %d", v, wid, seq, elapsed)
	return v, nil
}

func (f *SnowflakeFactory) BuildID2(wid, seq uint32, elapsed uint64) (uint64, error) {
	if elapsed >= 1<<f.bitsTimestamp {
		return 0, ErrOverTimeLimit
	}
	v := uint64(wid)<<(f.bitsSequence+f.bitsTimestamp) |
		elapsed<<f.bitsSequence |
		uint64(seq)
	// log.Printf("build %d w: %d seq: %2d elapsed: %d ", v, wid, seq, elapsed)
	return v, nil
}

func (f *SnowflakeFactory) NewSnowflake(wid uint32) (*Snowflake, error) {
	if wid > f.maxWorkerID {
		wid = f.maxWorkerID
	}
	// log.Printf("worker: %d len_worker: %d len_seq: %d len_ts: %d", wid, f.bitsWorkerID, f.bitsSequence, f.bitsTimestamp)
	return &Snowflake{f: f, build: f.BuildID2, worker: wid, mtx: &sync.Mutex{}}, nil
}

func NewSnowflake(worker uint32) (*Snowflake, error) {
	start, _ := time.Parse(time.RFC3339, "2018-10-24T07:30:06Z")
	return NewSnowflakeFactory(10, 12, 1, start).NewSnowflake(worker)
}

type Snowflake struct {
	f        *SnowflakeFactory
	build    func(uint32, uint32, uint64) (uint64, error)
	worker   uint32
	elapsed  uint64
	sequence uint32
	mtx      *sync.Mutex
}

func (s *Snowflake) WorkerID() uint32 { return s.worker }

func (s *Snowflake) ID() (uint64, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	elapsed := s.f.Elapsed()
	if s.elapsed < elapsed {
		s.elapsed = elapsed
		s.sequence = genRandomSequence(9)
		return s.build(s.worker, s.sequence, s.elapsed)
	}

	if s.elapsed > elapsed {
		elapsed = s.f.Elapsed()
		if s.elapsed > elapsed {
			return 0, ErrInvalidSystemClock
		}
	}

	s.sequence = s.f.MaskSequence(s.sequence + 1)
	if s.sequence == 0 {
		s.elapsed = s.elapsed + 1
		time.Sleep(s.f.Sleep(time.Duration(s.elapsed - elapsed)))
	}

	return s.build(s.worker, s.sequence, s.elapsed)
}

func genRandomSequence(n int32) uint32 {
	return uint32(rand.New(rand.NewSource(time.Now().In(time.UTC).UnixNano())).Int31n(n))
}

func WorkerIDFromIP(ipv4 net.IP) uint32 {
	if ipv4 == nil {
		return 0
	}
	ip := ipv4.To4()
	return uint32(ip[2])<<8 + uint32(ip[3])
}

func WorkerIDFromLocalIP() uint32 {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = os.Getenv("HOSTNAME")
	}

	var ipv4 net.IP
	addrs, _ := net.LookupIP(hostname)
	for _, addr := range addrs {
		if ipv4 = addr.To4(); ipv4 != nil {
			break
		}
	}
	return WorkerIDFromIP(ipv4)
}
