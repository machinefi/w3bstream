package types

import "strings"

type CommaSplitStrings []string

func (s CommaSplitStrings) MarshalText() ([]byte, error) {
	return []byte(strings.Join(s, ",")), nil
}

func (s *CommaSplitStrings) UnmarshalText(text []byte) error {
	*s = strings.Split(string(text), ",")
	return nil
}
