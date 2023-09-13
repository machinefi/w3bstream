package stream

import (
	"sync"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_buf"
)

// streamed buffer support TCP connection

type buffer struct {
	buf    []byte
	size   int // cap(buf)
	r, w   int
	rs, ws int // reading/writing available size
	mu     *sync.Mutex
}

const gDefaultSize = 20 * 1024 * 1024 // 20K default size

func New(cap int) *buffer {
	size := gDefaultSize
	if cap != 0 {
		size = cap
	}
	return &buffer{
		buf:  make([]byte, size),
		size: size,
		r:    0,
		w:    0,
		rs:   0,
		ws:   size,
		mu:   &sync.Mutex{},
	}
}

func (s *buffer) afterR(n int) {
	s.r += n
	s.r %= s.size
	s.rs -= n
	s.ws += n
}

func (s *buffer) afterW(n int) {
	s.w += n
	s.w %= s.size
	s.rs += n
	s.ws -= n
}

func (s *buffer) read(n int) ([]byte, error) {
	if s.rs < n {
		return nil, sock_buf.EStreamBufferDataLack
	}

	var ret = make([]byte, n)

	if s.r > s.w {
		if s.size-s.r >= n {
			copy(ret, s.buf[s.r:s.r+n])
		} else {
			c1 := s.size - s.r
			c2 := n - c1
			copy(ret, s.buf[s.r:s.r+c1])
			copy(ret[c1:], s.buf[0:c2])
		}
	} else {
		copy(ret, s.buf[s.r:s.r+n])
	}
	return ret, nil
}

func (s *buffer) resize() {
	buf := make([]byte, s.size<<1)
	ori, _ := s.read(s.rs)
	copy(buf[0:], ori[0:])
	s.buf = buf
	s.size <<= 1
	s.r, s.w, s.rs, s.ws = 0, len(ori), len(ori), s.size-len(ori)
}

func (s *buffer) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rs
}

func (s *buffer) Cap() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.size
}

func (s *buffer) Read(out []byte) (int, error) {
	n := len(out)
	if n == 0 {
		return 0, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	ret, err := s.read(n)
	if err != nil {
		return 0, err
	}
	s.afterR(n)
	copy(out, ret)
	return n, err
}

func (s *buffer) Write(in []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := len(in)

	for s.ws < n {
		s.resize()
	}

	defer s.afterW(n)
	if s.r > s.w {
		copy(s.buf[s.w:s.w+n], in)
	} else {
		if s.size-s.w >= n {
			copy(s.buf[s.w:s.w+n], in)
		} else {
			c1 := s.size - s.w
			c2 := n - c1
			copy(s.buf[s.w:s.size], in[0:c1])
			copy(s.buf[0:c2], in[c1:])
		}
	}
	return len(in), nil
}

func (s *buffer) ResetAndWrite(in []byte) (err error) {
	_, err = s.Write(in)
	return
}

func (s *buffer) WriteByte(v byte) error {
	_, err := s.Write([]byte{v})
	return err
}

func (s *buffer) Probe(n int) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if n > s.rs {
		return nil, sock_buf.EStreamBufferDataLack
	}

	return s.read(n)
}

func (s *buffer) Shift(n int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.afterR(n)
	return nil
}

func (s *buffer) Bytes() []byte {
	ret, _ := s.Probe(s.rs)
	return ret
}

func (s *buffer) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.r, s.w, s.rs, s.ws = 0, 0, 0, s.size
}
