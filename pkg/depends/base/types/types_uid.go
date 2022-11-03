package types

import "strconv"

type UID uint64

func (uid UID) MarshalText() ([]byte, error) {
	return []byte(uid.String()), nil
}

func (uid *UID) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	v, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return err
	}
	*uid = UID(v)
	return nil
}

func (uid UID) String() string { return strconv.FormatUint(uid.Uint(), 10) }

func (uid UID) Uint() uint64 { return uint64(uid) }

func AsUID(v uint64) UID { return UID(v) }

type UIDs []UID
