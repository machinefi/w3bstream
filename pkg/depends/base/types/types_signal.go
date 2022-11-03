package types

import (
	"strconv"
	"syscall"
)

type Signal syscall.Signal

func (s Signal) Int() int { return int(s) }

func (s Signal) String() string { return syscall.Signal(s).String() }

func (s Signal) MarshalText() ([]byte, error) {
	return []byte(strconv.Itoa(int(s))), nil
}

func (s *Signal) UnmarshalText(data []byte) error {
	v, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	*s = Signal(v)
	return nil
}

func (s Signal) Error() string { return syscall.Signal(s).String() }
