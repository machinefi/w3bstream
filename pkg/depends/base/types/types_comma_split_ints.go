package types

import (
	"bytes"
	"strconv"
)

type CommaSplitInts []int64

func (s CommaSplitInts) MarshalText() (data []byte, err error) {
	for _, i := range s {
		data = strconv.AppendInt(data, i, 10)
		data = append(data, ',')
	}
	data = bytes.TrimRight(data, ",")
	return
}

func (s *CommaSplitInts) UnmarshalText(text []byte) error {
	*s = (*s)[0:0]
	parts := bytes.Split(text, []byte{','})
	for _, data := range parts {
		v, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		*s = append(*s, v)
	}
	return nil
}
