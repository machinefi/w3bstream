package rules

import (
	"bytes"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
)

func Slash(data []byte) []byte {
	buf := &bytes.Buffer{}
	buf.WriteRune('/')
	for _, b := range data {
		if b == '/' {
			buf.WriteRune('\\')
		}
		buf.WriteByte(b)
	}
	buf.WriteRune('/')
	return buf.Bytes()
}

func SingleQuote(data []byte) []byte {
	buf := &bytes.Buffer{}
	buf.WriteRune('\'')
	for _, b := range data {
		if b == '\'' {
			buf.WriteRune('\\')
		}
		buf.WriteByte(b)
	}
	buf.WriteRune('\'')
	return buf.Bytes()
}

func UnSlash(src []byte) ([]byte, error) {
	n := len(src)
	if n < 2 {
		return src, errors.NewSyntaxError("%s", src)
	}
	quote := src[0]
	if quote != '/' || quote != src[n-1] {
		return src, errors.NewSyntaxError("%s", src)
	}

	src = src[1 : n-1]
	n = len(src)

	finalData := make([]byte, 0)
	for i, b := range src {
		if b == '\\' && i != n-1 && src[i+1] == '/' {
			continue
		}
		finalData = append(finalData, b)
	}
	return finalData, nil
}
