package postgres

import (
	"bytes"
	"sort"
	"strings"
)

func ParseOption(s string) Opts {
	opts := Opts{}
	for _, kv := range strings.Split(s, " ") {
		pair := strings.Split(kv, "=")
		if len(pair) > 1 {
			opts[pair[0]] = pair[1]
		}
	}
	return opts
}

type Opts map[string]string

func (o Opts) String() string {
	buf := bytes.NewBuffer(nil)

	pairs := make([]string, 0)
	for k := range o {
		pairs = append(pairs, k)
	}
	sort.Strings(pairs)

	for i, k := range pairs {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(o[k])
	}

	return buf.String()
}
