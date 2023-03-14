package types_test

import (
	"bytes"
	"net/url"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/base/types"
)

func TestAddress(t *testing.T) {
	cases1 := map[string]struct {
		Addr   *Address
		Expect string
	}{
		"asset": {
			Addr: &Address{
				Grp: "avatar",
				Key: "filename",
				Ext: "png",
			},
			Expect: "asset://avatar/filename.png",
		},
		"url": {
			Addr: &Address{
				URL: "http://demo.com/avatar/filename.png",
			},
			Expect: "http://demo.com/avatar/filename.png",
		},
		"without_ext": {
			Addr: &Address{
				Grp: "avatar",
				Key: "filename",
			},
			Expect: "asset://avatar/filename",
		},
		"local_file": {
			Addr: &Address{
				URL: "file:///AbsPath/To/Your/Local/File.ext",
			},
			Expect: "file:///AbsPath/To/Your/Local/File.ext",
		},
	}

	for name, c := range cases1 {
		t.Run("string:"+name, func(t *testing.T) {
			NewWithT(t).Expect(c.Addr.String()).To(Equal(c.Expect))
		})
	}

	cases2 := map[string]string{
		"asset":       "asset://avatar/filename.png",
		"url":         "http://venusgroup.com/avatar/filename.png",
		"without_ext": "asset://avatar/filename",
		"local_file":  "file:///AbsPath/To/Your/Local/File.ext",
	}

	for name, str := range cases2 {
		t.Run("parse:"+name, func(t *testing.T) {
			addr, err := ParseAddress(str)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(addr.String()).To(Equal(str))
		})
	}
}

func TestBinary(t *testing.T) {
	bytes, err := Binary("917").MarshalText()
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(bytes).To(Equal([]byte("917")))
}

func TestDuration(t *testing.T) {
	durations := map[string]struct {
		t Span
		s string
	}{
		"duration.h":  {AsDuration(time.Hour), "1h"},
		"duration.m":  {AsDuration(time.Minute), "1m"},
		"duration.s":  {AsDuration(time.Second), "1s"},
		"duration.ms": {AsDuration(time.Millisecond), "1ms"},
		"duration.us": {AsDuration(time.Microsecond), "1us"},
		"duration.ns": {AsDuration(time.Nanosecond), "1ns"},
		"duration1":   {AsDuration(Spans(Hours(1), Minutes(2), Milliseconds(3), Nanoseconds(4)).Duration()), "1h2m3ms4ns"},
		"duration2":   {AsDuration(Spans(Hours(1), Seconds(2), Microseconds(3)).Duration()), "1h2s3us"},
	}

	for name, d := range durations {
		t.Run(name, func(t *testing.T) {
			data, err := d.t.MarshalText()
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(data).To(Equal([]byte(d.s)))

			du := Duration(0)
			err = du.UnmarshalText([]byte(d.s))
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(du.Duration()).To(Equal(d.t.Duration()))
		})
	}
	t.Run("hour", func(t *testing.T) {
		data, err := Hour(1).MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(data).To(Equal([]byte("1h")))

		du := Hour(0)
		err = du.UnmarshalText([]byte("1h"))
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(du.Duration()).To(Equal(Hour(1).Duration()))
	})

	t.Run("minute", func(t *testing.T) {
		data, err := Minute(1).MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(data).To(Equal([]byte("1m")))

		du := Minute(0)
		err = du.UnmarshalText([]byte("1m"))
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(du.Duration()).To(Equal(Minute(1).Duration()))
	})

	t.Run("second", func(t *testing.T) {
		data, err := Second(1).MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(data).To(Equal([]byte("1s")))

		du := Second(0)
		err = du.UnmarshalText([]byte("1s"))
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(du.Duration()).To(Equal(Second(1).Duration()))
	})

	t.Run("millisecond", func(t *testing.T) {
		data, err := Millisecond(1).MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(data).To(Equal([]byte("1ms")))

		du := Millisecond(0)
		err = du.UnmarshalText([]byte("1ms"))
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(du.Duration()).To(Equal(Millisecond(1).Duration()))
	})

	t.Run("microsecond", func(t *testing.T) {
		data, err := Microsecond(1).MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(data).To(Equal([]byte("1us")))

		du := Microsecond(0)
		err = du.UnmarshalText([]byte("1us"))
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(du.Duration()).To(Equal(Microsecond(1).Duration()))
	})

	t.Run("nanosecond", func(t *testing.T) {
		data, err := Nanosecond(1).MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(data).To(Equal([]byte("1ns")))

		du := Nanosecond(0)
		err = du.UnmarshalText([]byte("1ns"))
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(du.Duration()).To(Equal(Nanosecond(1).Duration()))
	})
}

func TestEndpoint(t *testing.T) {
	cases := map[string]struct {
		ep  *Endpoint
		url string
	}{
		"STMPs": {
			ep: &Endpoint{
				Scheme:   "stmps",
				Hostname: "stmps.mail.xxx.com",
				Port:     465,
			},
			url: "stmps://stmps.mail.xxx.com:465",
		},
		"Postgres": {
			ep: &Endpoint{
				Scheme:   "postgres",
				Hostname: "hostname",
				Username: "username",
				Password: "password",
				Port:     5432,
				Base:     "database_name",
				Param:    url.Values{"sslmode": {"disable"}},
			},
			url: "postgres://username:password@hostname:5432/database_name?sslmode=disable",
		},
		"NoScheme": {
			ep: &Endpoint{
				Scheme:   "",
				Hostname: "hostname",
				Username: "username",
				Password: "password",
				Port:     5432,
				Base:     "database_name",
				Param:    url.Values{"sslmode": {"disable"}},
			},
			url: "//username:password@hostname:5432/database_name?sslmode=disable",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if c.ep != nil {
				NewWithT(t).Expect(c.ep.String()).To(Equal(c.url))
			}
			parsed, err := ParseEndpoint(c.url)
			if c.ep != nil {
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(*parsed).To(Equal(*c.ep))
			} else {
				NewWithT(t).Expect(err).NotTo(BeNil())
			}
		})
	}

	t.Run("UnmarshalExtra", func(t *testing.T) {
		opt := struct {
			ConnectTimeout Duration `name:"connectTimeout" default:"10s"`
			ReadTimeout    Duration `name:"readTimeout" default:"10s"`
			WriteTimeout   Duration `name:"writeTimeout" default:"10s"`
			IdleTimeout    Duration `name:"idleTimeout" default:"240s"`
			MaxActive      int      `name:"maxActive" default:"5"`
			MaxIdle        int      `name:"maxIdle" default:"3"`
			DB             int      `name:"db" default:"10"`
		}{}

		err := UnmarshalExtra(url.Values{}, &opt)
		NewWithT(t).Expect(err).To(BeNil())
		// spew.Dump(opt)
	})
}

func TestSignal(t *testing.T) {
}

func TestEthAddress(t *testing.T) {
	ori := []byte("0x3110752149AF23Ee65968C2019b7c86D12B32229")
	addr := EthAddress("")

	err := addr.UnmarshalText(ori)
	NewWithT(t).Expect(err).To(BeNil())

	cmp, err := addr.MarshalText()
	NewWithT(t).Expect(cmp).To(Equal(bytes.ToLower(ori)))
}
