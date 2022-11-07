package types

import (
	"strconv"
	"strings"
)

type Capacity int64

type MB float64

func (v MB) MarshalText() ([]byte, error) {
	return strconv.AppendFloat([]byte{}, float64(v), 'g', -1, 64), nil
}

func (v *MB) UnmarshalText(data []byte) error {
	str := strings.ToUpper(string(data))
	mb, err := strconv.ParseFloat(strings.TrimRight(str, "MB"), 64)
	if err != nil {
		return err
	}
	*v = MB(mb)
	return nil
}

type KB int64

type GB int64
