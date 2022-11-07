package validator

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type StrLenMode int

const (
	STR_LEN_MODE__LENGTH StrLenMode = iota
	STR_LEN_MODE__RUNE_COUNT
)

var StrLenModes = map[StrLenMode]func(s string) uint64{
	STR_LEN_MODE__LENGTH:     func(s string) uint64 { return uint64(len(s)) },
	STR_LEN_MODE__RUNE_COUNT: func(s string) uint64 { return uint64(utf8.RuneCount([]byte(s))) },
}

func ParseStrLenMode(s string) (StrLenMode, error) {
	switch strings.ToLower(s) {
	case "rune_count":
		return STR_LEN_MODE__RUNE_COUNT, nil
	case "length", "":
		return STR_LEN_MODE__LENGTH, nil
	default:
		return STR_LEN_MODE__LENGTH, fmt.Errorf("unsupported string length mode")
	}
}

func (m StrLenMode) String() string {
	switch m {
	case STR_LEN_MODE__RUNE_COUNT:
		return "rune_count"
	default:
		return "length"
	}
}

func StringRuleParam(r *Rule) (mode StrLenMode, err error) {
	if r.Params != nil {
		if len(r.Params) != 1 {
			err = fmt.Errorf(
				"string should only 1 parameter, but got %d",
				len(r.Params),
			)
			return
		}
		mode, err = ParseStrLenMode(string(r.Params[0].Bytes()))
		return
	}
	if r.Name == "char" {
		return STR_LEN_MODE__RUNE_COUNT, nil
	}
	return STR_LEN_MODE__LENGTH, nil
}

func StringRuleValues(r *Rule) map[string]string {
	values := r.ComputedValues()
	if values == nil {
		return nil
	}
	enums := map[string]string{}
	for _, v := range values {
		str := string(v.Bytes())
		enums[str] = str
	}
	return enums
}
