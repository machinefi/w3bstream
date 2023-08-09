package logger

import (
	"go.opentelemetry.io/otel/attribute"

	"github.com/machinefi/w3bstream/pkg/depends/x/textx"
)

func KVsToAttr(kvs ...any) (atts []attribute.KeyValue) {
	n := len(kvs)
	if n > 0 && n%2 == 0 {
		atts = make([]attribute.KeyValue, n/2)
		for i := range atts {
			k, v := kvs[2*i], kvs[2*i+1]

			if key, ok := k.(string); ok {
				val, err := textx.MarshalText(v)
				if err != nil {
					continue
				}
				atts[i].Key = attribute.Key(key)
				atts[i].Value = attribute.StringValue(string(val))
			}
		}
		return atts
	}
	return nil
}
