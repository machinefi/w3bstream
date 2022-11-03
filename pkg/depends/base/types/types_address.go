package types

import (
	"database/sql/driver"
	"fmt"
	"net/url"
	"strings"
)

type Address struct {
	URL string
	Grp string
	Key string
	Ext string
}

var (
	_ String          = (*Address)(nil)
	_ TextMarshaler   = (*Address)(nil)
	_ TextUnmarshaler = (*Address)(nil)
)

const AddressSchemeName = "asset"

func ParseAddress(text string) (*Address, error) {
	u, err := url.Parse(text)
	if err != nil {
		return nil, err
	}
	a := &Address{}

	if u.Scheme != AddressSchemeName {
		a.URL = u.String()
		return a, nil
	}
	a.Grp = u.Hostname()
	if len(u.Path) > 0 {
		a.Key = u.Path[1:]
		if idx := strings.LastIndex(u.Path, "."); idx != -1 {
			a.Key = u.Path[1:idx]
			a.Ext = u.Path[idx+1:]
		}
	}
	return a, nil
}

func (a Address) String() string {
	if a.URL != "" {
		return a.URL
	}
	if a.Grp == "" && a.Key == "" {
		return ""
	}
	u := fmt.Sprintf("%s://%s/%s", AddressSchemeName, a.Grp, a.Key)
	if a.Ext != "" {
		u += "." + a.Ext
	}
	return u
}

func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *Address) UnmarshalText(text []byte) error {
	addr, err := ParseAddress(string(text))
	if err != nil {
		return err
	}
	*a = *addr
	return nil
}

func (a Address) DataType(string) string { return "varchar(1024)" }

func (a Address) Value() (driver.Value, error) { return a.String(), nil }

func (a *Address) Scan(src interface{}) error {
	return a.UnmarshalText([]byte(src.(string)))
}
