package env_test

import (
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	. "github.com/machinefi/w3bstream/pkg/depends/conf/env"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
)

func TestPathWalker(t *testing.T) {
	pw := NewPathWalker()

	pw.Enter("parent")
	NewWithT(t).Expect(pw.Paths()).To(Equal([]interface{}{"parent"}))
	NewWithT(t).Expect(pw.String()).To(Equal("parent"))

	pw.Enter("child")
	NewWithT(t).Expect(pw.Paths()).To(Equal([]interface{}{"parent", "child"}))
	NewWithT(t).Expect(pw.String()).To(Equal("parent_child"))

	pw.Enter("prop")
	NewWithT(t).Expect(pw.Paths()).To(Equal([]interface{}{"parent", "child", "prop"}))
	NewWithT(t).Expect(pw.String()).To(Equal("parent_child_prop"))

	pw.Enter(100)
	NewWithT(t).Expect(pw.Paths()).To(Equal([]interface{}{"parent", "child", "prop", 100}))
	NewWithT(t).Expect(pw.String()).To(Equal("parent_child_prop_100"))

	pw.Exit()
	NewWithT(t).Expect(pw.Paths()).To(Equal([]interface{}{"parent", "child", "prop"}))
	NewWithT(t).Expect(pw.String()).To(Equal("parent_child_prop"))

	pw.Exit()
	NewWithT(t).Expect(pw.Paths()).To(Equal([]interface{}{"parent", "child"}))
	NewWithT(t).Expect(pw.String()).To(Equal("parent_child"))
}

type SubConfig struct {
	Duration     types.Duration
	Password     types.Password `env:""`
	Key          string         `env:""`
	Bool         bool
	Map          map[string]string
	Func         func() error
	defaultValue bool
}

func (c *SubConfig) SetDefault() {
	c.defaultValue = true
}

type Config struct {
	Map       map[string]string
	Slice     []string `env:""`
	PtrString *string  `env:""`
	Host      string   `env:""`
	SubConfig
	Config SubConfig
}

func TestEnvVars(t *testing.T) {
	c := Config{}

	c.Duration = types.Duration(time.Hour + 2*time.Second + 3*time.Microsecond)
	c.Password = types.Password("123123")
	c.Key = "123456"
	c.PtrString = ptrx.String("123456=")
	c.Slice = []string{"1", "2"}
	c.Config.Key = "key"
	c.Config.defaultValue = true
	c.defaultValue = true

	envVars := NewVars("S")

	t.Run("Encoding", func(t *testing.T) {
		data, _ := NewEncoder(envVars).Encode(&c)
		NewWithT(t).Expect(string(data)).To(
			Equal(`S__Bool=false
S__Config_Bool=false
S__Config_Duration=0s
S__Config_Key=key
S__Config_Password=
S__Duration=1h2s3us
S__Host=
S__Key=123456
S__Password=123123
S__PtrString=123456=
S__Slice_0=1
S__Slice_1=2
`))
	})

	t.Run("Decoding", func(t *testing.T) {
		data, _ := NewEncoder(envVars).Encode(&c)

		envVars := LoadVarsFromEnviron("S", strings.Split(string(data), "\n"))

		c2 := Config{}
		err := NewDecoder(envVars).Decode(&c2)

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(c2).To(Equal(c))
	})
}
