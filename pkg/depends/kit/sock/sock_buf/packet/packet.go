package packet

import (
	"bytes"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_buf"
)

var _ = bytes.Buffer{}

var gBufferSize = 2 * 1024 * 1024

type buffer struct {
	*bytes.Buffer
}

func New(cap int) *buffer {
	size := gBufferSize
	if cap != 0 {
		size = cap
	}
	ret := &buffer{bytes.NewBuffer(make([]byte, size))}
	ret.Reset()
	return ret
}

func (b *buffer) Read(out []byte) (int, error) {
	return b.Buffer.Read(out)
}

func (b *buffer) Write(in []byte) (int, error) {
	return b.Buffer.Write(in)
}

func (b *buffer) ResetAndWrite(in []byte) error {
	b.Reset()

	var size int
	for size < len(in) {
		l, e := b.Write(in[size:])
		if e != nil {
			return e
		}
		size += l
	}
	return nil
}

func (b *buffer) Probe(n int) ([]byte, error) {
	if b.Len() < n {
		return nil, sock_buf.EPacketBufferDataLack
	}
	return b.Bytes()[0:n], nil
}

func (b *buffer) Shift(n int) error {
	_ = b.Buffer.Next(n)
	return nil
}

func (b *buffer) Bytes() []byte {
	return b.Buffer.Bytes()
}

func (b *buffer) Reset() {
	b.Buffer.Reset()
}
