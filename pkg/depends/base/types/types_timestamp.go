package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"strconv"
	"time"
)

var (
	UTC               = time.UTC
	CST               = time.FixedZone("CST", 8*60*60)
	TimestampZero     = Timestamp{time.Time{}}
	TimestampUnixZero = Timestamp{time.Unix(0, 0)}
)

const (
	DatetimeFormat = time.RFC3339
)

// openapi:strfmt date-time
type Timestamp struct{ time.Time }

var (
	_ sql.Scanner              = (*Timestamp)(nil)
	_ driver.Valuer            = (*Timestamp)(nil)
	_ encoding.TextMarshaler   = (*Timestamp)(nil)
	_ encoding.TextUnmarshaler = (*Timestamp)(nil)
)

func AsTimestamp(t time.Time) Timestamp { return Timestamp{t} }

func (Timestamp) DataType(engine string) string {
	return "bigint"
}

func ParseTimestampFromString(s string) (Timestamp, error) {
	var (
		t   Timestamp
		err error
	)
	t.Time, err = time.Parse(DatetimeFormat, s)
	if err != nil {
		return Timestamp{}, err
	}
	return t, nil
}

func ParseTimestampFromStringWithLayout(input, layout string) (Timestamp, error) {
	t, err := time.ParseInLocation(layout, input, CST)
	if err != nil {
		return TimestampUnixZero, err
	}
	return Timestamp{t}, nil
}

func (t *Timestamp) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return fmt.Errorf("sql.Scan() strfmt.Timestamp from: %#v failed: %s", v, err.Error())
		}
		*t = Timestamp{time.Unix(n, 0)}
	case int64:
		if v < 0 {
			*t = Timestamp{}
		} else {
			*t = Timestamp{time.Unix(v, 0)}
		}
	case nil:
		*t = TimestampZero
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Timestamp from: %#v", v)
	}
	return nil
}

func (t Timestamp) Value() (driver.Value, error) {
	s := t.Unix()
	if s < 0 {
		s = 0
	}
	return s, nil
}

func (t *Timestamp) Set(_t time.Time) { t.Time = _t }

func (t Timestamp) String() string {
	if t.IsZero() {
		return ""
	}
	return t.In(CST).Format(time.RFC3339)
}

func (t Timestamp) Format(layout string) string {
	return t.In(CST).Time.Format(layout)
}

func (t Timestamp) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Timestamp) UnmarshalText(data []byte) (err error) {
	str := string(data)
	if len(str) == 0 || str == "0" {
		return nil
	}
	*t, err = ParseTimestampFromString(str)
	return
}

func (t Timestamp) Unix() int64 { return t.Time.Unix() }

func (t Timestamp) IsZero() bool {
	unix := t.Unix()
	return unix == 0 || unix == TimestampZero.Unix()
}

func (t Timestamp) In(loc *time.Location) Timestamp {
	return Timestamp{t.Time.In(loc)}
}
