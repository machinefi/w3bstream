package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
)

type Duration time.Duration

func AsDuration(d time.Duration) *Duration { return (*Duration)(ptrx.Duration(d)) }

func (d *Duration) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	du, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = Duration(du)
	return nil
}

func (d Duration) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

func (d Duration) Duration() time.Duration { return time.Duration(d) }

func (d Duration) Int() int64 { return int64(d) }

func (d Duration) String() string {
	var (
		du    = int64(d)
		neg   = du < 0
		sep   [6]int64
		units = [6]string{"ns", "us", "ms", "s", "m", "h"}
		ret   = ""
	)
	if du == 0 {
		return "0s"
	}
	if neg {
		du = 0 - du
	}
	for i := 0; i < 6; i++ {
		if i < 3 {
			sep[i] = du % 1000
			du /= 1000
		} else if i < 5 {
			sep[i] = du % 60
			du /= 60
		} else {
			sep[i] = du
		}

	}
	if neg {
		ret += "-"
	}
	first, last := 5, 0
	for first >= 0 && sep[first] == 0 {
		first--
	}
	for last <= 5 && sep[last] == 0 {
		last++
	}
	for i := first; i >= last; i-- {
		if sep[i] == 0 {
			continue
		}
		ret += fmt.Sprintf("%d%s", sep[i], units[i])
	}
	return ret
}

func (d Duration) Literal() string { return strconv.FormatInt(d.Int(), 10) }

func (d Duration) Span() Span { return &d }

type Hour int64

func Hours(h int64) *Hour { return (*Hour)(ptrx.Int64(h)) }

func (h *Hour) UnmarshalText(text []byte) error {
	d := Duration(0)
	if err := d.UnmarshalText(text); err != nil {
		return err
	}
	*h = Hour((time.Duration)(d).Hours())
	return nil
}

func (h Hour) MarshalText() ([]byte, error) { return []byte(h.String()), nil }

func (h Hour) Duration() time.Duration { return time.Duration(h) * time.Hour }

func (h Hour) Int() int64 { return int64(h) }

func (h Hour) String() string { return h.Literal() + "h" }

func (h Hour) Literal() string { return strconv.FormatFloat(float64(h), 'f', -1, 64) }

func (h Hour) Span() Span { return &h }

type Minute int64

func Minutes(h int64) *Minute { return (*Minute)(ptrx.Int64(h)) }

func (m *Minute) UnmarshalText(text []byte) error {
	d := Duration(0)
	if err := d.UnmarshalText(text); err != nil {
		return err
	}
	*m = Minute((time.Duration)(d).Minutes())
	return nil
}

func (m Minute) MarshalText() ([]byte, error) { return []byte(m.String()), nil }

func (m Minute) Duration() time.Duration { return time.Duration(m) * time.Minute }

func (m Minute) Int() int64 { return int64(m) }

func (m Minute) String() string { return m.Literal() + "m" }

func (m Minute) Literal() string { return strconv.FormatInt(int64(m), 10) }

func (m Minute) Span() Span { return &m }

type Second int64

func Seconds(s int64) *Second { return (*Second)(ptrx.Int64(s)) }

func (s *Second) UnmarshalText(text []byte) error {
	str := strings.TrimRight(strings.ToLower(string(text)), "s")
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*s = Second(v)
	return nil
}

func (s Second) MarshalText() ([]byte, error) { return []byte(s.String()), nil }

func (s Second) Duration() time.Duration { return time.Duration(s) * time.Second }

func (s Second) Int() int64 { return int64(s) }

func (s Second) String() string { return s.Literal() + "s" }

func (s Second) Literal() string { return strconv.FormatInt(int64(s), 10) }

func (s Second) Span() Span { return &s }

type Millisecond int64

func Milliseconds(s int64) *Millisecond { return (*Millisecond)(ptrx.Int64(s)) }

func (ms *Millisecond) UnmarshalText(text []byte) error {
	d := Duration(0)
	if err := d.UnmarshalText(text); err != nil {
		return err
	}
	*ms = Millisecond((time.Duration)(d).Milliseconds())
	return nil
}

func (ms Millisecond) MarshalText() ([]byte, error) { return []byte(ms.String()), nil }

func (ms Millisecond) Duration() time.Duration { return time.Duration(ms) * time.Millisecond }

func (ms Millisecond) Int() int64 { return int64(ms) }

func (ms Millisecond) String() string { return ms.Literal() + "ms" }

func (ms Millisecond) Literal() string { return strconv.FormatInt(int64(ms), 10) }

func (ms Millisecond) Span() Span { return &ms }

type Microsecond int64

func Microseconds(v int64) *Microsecond { return (*Microsecond)(ptrx.Int64(v)) }

func (us *Microsecond) UnmarshalText(text []byte) error {
	d := Duration(0)
	if err := d.UnmarshalText(text); err != nil {
		return err
	}
	*us = Microsecond((time.Duration)(d).Microseconds())
	return nil
}

func (us Microsecond) MarshalText() ([]byte, error) { return []byte(us.String()), nil }

func (us Microsecond) Duration() time.Duration { return time.Duration(us) * time.Microsecond }

func (us Microsecond) Int() int64 { return int64(us) }

func (us Microsecond) String() string { return us.Literal() + "us" }

func (us Microsecond) Literal() string { return strconv.FormatInt(int64(us), 10) }

func (us Microsecond) Span() Span { return &us }

type Nanosecond int64

func Nanoseconds(v int64) *Nanosecond { return (*Nanosecond)(ptrx.Int64(v)) }

func (ns *Nanosecond) UnmarshalText(text []byte) error {
	d := Duration(0)
	if err := d.UnmarshalText(text); err != nil {
		return err
	}
	*ns = Nanosecond((time.Duration)(d).Nanoseconds())
	return nil
}

func (ns Nanosecond) MarshalText() ([]byte, error) { return []byte(ns.String()), nil }

func (ns Nanosecond) Duration() time.Duration { return time.Duration(ns) * time.Nanosecond }

func (ns Nanosecond) Int() int64 { return int64(ns) }

func (ns Nanosecond) String() string { return ns.Literal() + "ns" }

func (ns Nanosecond) Literal() string { return strconv.FormatInt(int64(ns), 10) }

func (ns Nanosecond) Span() Span { return &ns }

func Spans(du ...Span) Duration {
	sum := time.Duration(0)
	for _, span := range du {
		sum += span.Duration()
	}
	return Duration(sum)
}

var (
	_ Span = Hours(0)
	_ Span = Minutes(0)
	_ Span = Seconds(0)
	_ Span = Milliseconds(0)
	_ Span = Microseconds(0)
	_ Span = Nanoseconds(0)
	_ Span = AsDuration(0)
)
