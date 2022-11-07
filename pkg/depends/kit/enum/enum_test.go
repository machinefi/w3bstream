package enum_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

func TestScanEnum(t *testing.T) {
	cases := []struct {
		offset int
		values []interface{}
		expect []int
	}{
		{
			-3,
			[]interface{}{
				nil,
				[]byte("-3"),
				"-2",
				int(-1),
				int8(0),
				int16(1),
				int32(2),
				int64(3),
				uint(4),
				uint8(5),
				uint16(6),
				uint32(7),
				uint64(8),
			},
			[]int{
				0,
				0,
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8,
				9,
				10,
				11,
			},
		},
	}

	for _, c := range cases {
		for i, v := range c.values {
			n, err := enum.ScanIntEnumStringer(v, c.offset)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(n).To(Equal(c.expect[i]))
		}
	}
}
