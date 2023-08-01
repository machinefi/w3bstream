package types

import (
	"fmt"
	"go/ast"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/textx"
)

type Endpoint struct {
	Scheme   string
	Hostname string
	Port     uint16
	Base     string
	Path     string
	Username string
	Password Password
	Param    url.Values
}

func (e Endpoint) String() string {
	u := url.URL{
		Scheme:   e.Scheme,
		Host:     e.Host(),
		RawPath:  "",
		RawQuery: e.Param.Encode(),
	}

	if e.Base != "" {
		u.Path = "/" + e.Base
	}

	if e.Param != nil {
		u.RawQuery = e.Param.Encode()
	}

	if e.Username != "" || e.Password != "" {
		u.User = url.UserPassword(e.Username, e.Password.String())
	}

	s, err := url.QueryUnescape(u.String())
	if err != nil {
		log.Fatal(err)
	}
	return s
}

// URL with host, port and path only
func (e Endpoint) URL() string {
	return (&url.URL{
		Scheme: e.Scheme,
		Host:   e.Host(),
		Path:   e.Path,
	}).String()
}

func (e Endpoint) SecurityString() string {
	e.Password = Password(e.Password.SecurityString())
	return e.String()
}

func (e Endpoint) IsZero() bool { return e.Hostname == "" }

func (e Endpoint) IsTLS() bool {
	return !e.IsZero() && e.Scheme[len(e.Scheme)-1] == 's'
}

func (e Endpoint) Host() string {
	if e.Port == 0 {
		return e.Hostname
	}
	return e.Hostname + ":" + strconv.FormatUint(uint64(e.Port), 10)
}

func (e Endpoint) SchemeHost() string {
	host := e.Host()
	if e.Scheme == "" {
		return host
	}
	return fmt.Sprintf("%s://%s", e.Scheme, host)
}

func (e Endpoint) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *Endpoint) UnmarshalText(text []byte) error {
	ep, err := ParseEndpoint(string(text))
	if err != nil {
		return err
	}
	*e = *ep
	return nil
}

func ParseEndpoint(text string) (*Endpoint, error) {
	u, err := url.Parse(text)
	if err != nil {
		return nil, err
	}

	ep := &Endpoint{}

	ep.Scheme = u.Scheme

	if q := u.Query(); len(q) > 0 {
		ep.Param = q
	}

	ep.Path = u.Path

	if len(u.Path) > 0 {
		ep.Base = strings.Split(u.Path[1:], "/")[0]
	}

	ep.Hostname = u.Hostname()

	if port, err := strconv.ParseUint(u.Port(), 10, 16); err == nil {
		ep.Port = uint16(port)
	}

	if u.User != nil {
		ep.Username = u.User.Username()
		password, _ := u.User.Password()
		ep.Password = Password(password)
	}

	return ep, nil
}

func UnmarshalExtra(ext url.Values, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointor value %s is not supported", rv.Type())
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return nil
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		ft := rt.Field(i)
		fn := ft.Name

		if !ast.IsExported(fn) {
			continue
		}
		if tag, ok := ft.Tag.Lookup("name"); ok {
			n, _ := reflectx.TagValueAndFlags(tag)
			if n == "-" {
				continue
			}
			if n != "" {
				fn = n
			}
		}
		fv := rv.Field(i)
		value := ext.Get(fn)
		if value == "" {
			value = ft.Tag.Get("default")
		}
		if err := textx.UnmarshalText(fv, []byte(value)); err != nil {
			return err
		}
	}
	return nil
}
