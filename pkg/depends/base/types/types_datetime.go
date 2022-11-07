package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"time"
)

var (
	DatetimeZero = Datetime{time.Time{}}
)

// openapi:strfmt date-time
type Datetime struct{ time.Time }

var (
	_ sql.Scanner              = (*Datetime)(nil)
	_ driver.Valuer            = (*Datetime)(nil)
	_ encoding.TextMarshaler   = (*Datetime)(nil)
	_ encoding.TextUnmarshaler = (*Datetime)(nil)
)

func AsDatetime(t time.Time) Datetime { return Datetime{t} }

func (Datetime) DataType(engine string) string { return "timestamp" }

func ParseDatetimeFromString(s string) (Datetime, error) {
	t, err := time.Parse(DatetimeFormat, s)
	if err != nil {
		return Datetime{}, err
	}
	return Datetime{t}, nil
}

func ParseDatetimeFromStringWithFormatterInCST(s, formatter string) (Datetime, error) {
	t, err := time.ParseInLocation(formatter, s, CST)
	if err != nil {
		return Datetime{}, err
	}
	return Datetime{t}, nil
}

func (t *Datetime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		t.Time = time.Unix(v.Unix(), 0)
	case nil:
		*t = DatetimeZero
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Datetime from: %#v", v)
	}
	return nil
}

func (t *Datetime) Set(_t time.Time) { t.Time = _t }

func (t Datetime) Value() (driver.Value, error) {
	return time.Unix(t.Unix(), 0), nil
}

func (t Datetime) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Time.In(CST).Format(DatetimeFormat)
}

func (t Datetime) Format(layout string) string {
	return t.Time.In(CST).Format(layout)
}

func (t Datetime) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Datetime) UnmarshalText(data []byte) (err error) {
	str := string(data)
	if len(str) == 0 || str == "0" {
		return nil
	}
	*t, err = ParseDatetimeFromString(str)
	return
}

func (t Datetime) IsZero() bool {
	unix := t.Unix()
	return unix == 0 || unix == DatetimeZero.Unix()
}

func (t Datetime) In(loc *time.Location) Datetime {
	return Datetime{t.Time.In(loc)}
}
