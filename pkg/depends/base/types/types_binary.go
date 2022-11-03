package types

type Binary []byte

var (
	_ TextMarshaler   = (Binary)(nil)
	_ TextUnmarshaler = (*Binary)(nil)
)

func (d Binary) MarshalText() ([]byte, error) { return d, nil }

func (d *Binary) UnmarshalText(data []byte) (err error) {
	*d = Binary(data)
	return
}
