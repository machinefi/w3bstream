package types

import "strconv"

// Number.MAX_SAFE_INTEGER JavaScript (2^53 – 1)

type SFID uint64

func (sf SFID) MarshalText() ([]byte, error) { return []byte(sf.String()), nil }

func (sf *SFID) UnmarshalText(data []byte) (err error) {
	str := string(data)
	if len(str) == 0 {
		return
	}
	var u uint64
	u, err = strconv.ParseUint(str, 10, 64)
	*sf = SFID(u)
	return
}

func (sf SFID) String() string { return strconv.FormatUint(uint64(sf), 10) }

type SFIDs []SFID

func (sfs SFIDs) ToUint64() []uint64 {
	var l []uint64
	for _, sf := range sfs {
		l = append(l, uint64(sf))
	}
	return l
}
