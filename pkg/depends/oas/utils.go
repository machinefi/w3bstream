package oas

import (
	"bytes"
	"encoding/json"
	"log"
)

var (
	comma   = byte(',')
	closers = map[byte]byte{
		'{': '}',
		'[': ']',
	}
)

func concatJSON(blobs ...[]byte) []byte {
	if len(blobs) == 0 {
		return nil
	}
	if len(blobs) == 1 {
		return blobs[0]
	}

	last := len(blobs) - 1
	var opening, closing byte
	a := 0
	idx := 0
	buf := bytes.NewBuffer(nil)

	for i, b := range blobs {
		if len(b) > 0 && opening == 0 {
			opening, closing = b[0], closers[b[0]]
		}

		if opening != '{' && opening != '[' {
			continue
		}

		if len(b) < 3 {
			if i == last && a > 0 {
				if err := buf.WriteByte(closing); err != nil {
					log.Println(err)
				}
			}
			continue
		}

		idx = 0
		if a > 0 {
			if err := buf.WriteByte(comma); err != nil {
				log.Println(err)
			}
			idx = 1
		}

		if i != last {
			if _, err := buf.Write(b[idx : len(b)-1]); err != nil {
				log.Println(err)
			}
		} else {
			if _, err := buf.Write(b[idx:]); err != nil {
				log.Println(err)
			}
		}
		a++
	}

	if buf.Len() == 0 {
		if err := buf.WriteByte(opening); err != nil {
			log.Println(err)
		}
		if err := buf.WriteByte(closing); err != nil {
			log.Println(err)
		}
	}
	return buf.Bytes()
}

func flattenMarshalJSON(values ...interface{}) ([]byte, error) {
	blobs := make([][]byte, 0)
	for i := range values {
		v := values[i]
		if v != nil {
			b, err := json.Marshal(v)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			blobs = append(blobs, b)
		}
	}
	return concatJSON(blobs...), nil
}

func flattenUnmarshalJSON(data []byte, values ...interface{}) error {
	for i := range values {
		err := json.Unmarshal(data, values[i])
		if err != nil {
			return err
		}
	}
	return nil
}
