package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/deploy"
	"github.com/machinefi/w3bstream/pkg/depends/conf/env"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

type Ctx struct {
	cmd       *cobra.Command
	name      string          // name app name
	feat      string          // feat git feature
	version   string          // version git version|git tag
	root      string          // root app root
	vars      []*env.Vars     // vars default env vars
	conf      []reflect.Value // conf config reflect.Value
	deployers map[string]deploy.Deployer
	ctx       context.Context
}

func New(setters ...OptSetter) *Ctx {
	c := &Ctx{ctx: context.Background()}
	for _, setter := range setters {
		setter(c)
	}
	c.cmd = &cobra.Command{}
	if feat, ok := os.LookupEnv(consts.EnvProjectFeat); ok && feat != "" {
		c.feat = feat
	}
	if version, ok := os.LookupEnv(consts.EnvProjectVersion); ok && version != "" {
		c.version = version
	}
	if name, ok := os.LookupEnv(consts.EnvProjectName); ok && name != "" {
		c.name = name
	}
	_ = os.Setenv(consts.EnvProjectName, c.name)
	return c
}

func (c *Ctx) Context() context.Context { return c.ctx }

// Conf init all configs from yml file, and do initialization for each config.
// config dir include `config.yml.template` `config.yml` and `master.yml`
// config.yml.template shows config file template and preset config values
// config.yml contains all user configured values
func (c *Ctx) Conf(configs ...interface{}) {
	local, err := os.ReadFile(filepath.Join(c.root, "./config/config.yml"))
	if err == nil {
		kv := make(map[string]string)
		if err = yaml.Unmarshal(local, &kv); err == nil {
			for k, v := range kv {
				_ = os.Setenv(k, v)
			}
		}
	}

	if key := os.Getenv(consts.GoRuntimeEnv); key == "" {
		_ = os.Setenv(consts.GoRuntimeEnv, consts.ProduceEnv)
	}

	for _, v := range configs {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr {
			panic("should pass pointer for setting value")
		}

		must.NoError(c.scan(rv))
		must.NoError(c.marshal(rv))
	}

	if err = c.MarshalDefault(); err != nil {
		panic(err)
	}

	for _, v := range configs {
		rv := reflect.ValueOf(v)
		c.conf = append(c.conf, rv)

		if zero, ok := v.(types.ZeroChecker); ok && zero.IsZero() {
			t := reflect.Indirect(reflect.ValueOf(v)).Type()
			log.Println(errors.Errorf(
				"zero config: %s",
				color.CyanString("[%s.%s]:", filepath.Base(t.PkgPath()), t.Name()),
			))
		}

		switch conf := v.(type) {
		case interface{ Init() }:
			conf.Init()
		case interface{ Init() error }:
			if err = conf.Init(); err != nil {
				panic(errors.Errorf("conf init: %v", err))
			}
		}

		rv = reflectx.Indirect(rv)
		if rv.Kind() == reflect.Struct {
			for i := 0; i < rv.NumField(); i++ {
				value := rv.Field(i)
				if !value.CanInterface() {
					continue
				}
				if value.Type().Kind() == reflect.Interface {
					panic("interface type unsupported in config scanning")
				}
				fv := value.Interface()
				ft := reflect.Indirect(reflect.ValueOf(fv)).Type()
				if zero, ok := fv.(types.ZeroChecker); ok && zero.IsZero() {
					log.Println(errors.Errorf(
						"zero config: %s",
						color.CyanString("[%s.%s]:", filepath.Base(ft.PkgPath()), ft.Name()),
					))
					continue
				}
				switch conf := value.Interface().(type) {
				case interface{ Init() }:
					conf.Init()
				case interface{ Init() error }:
					if err = conf.Init(); err != nil {
						panic(errors.Errorf("init failed %s %s",
							color.CyanString("[%s.%s]:", filepath.Base(ft.PkgPath()), ft.Name()),
							color.RedString("[%v]", err),
						))
					}
				}
			}
		}
	}
}

func (c *Ctx) AddCommand(name string, fn func(...string), commands ...func(*cobra.Command)) {
	cmd := &cobra.Command{Use: name}

	for i := range commands {
		commands[i](cmd)
	}

	cmd.Run = func(_ *cobra.Command, args []string) {
		fn(args...)
	}

	c.cmd.AddCommand(cmd)
}

func (c *Ctx) String() string {
	ret := c.name
	if c.feat != "" {
		ret += "--" + c.feat
	}
	if c.version != "" {
		ret += "@" + c.version
	}
	return ret
}

func (c *Ctx) Root() string { return c.root }

func (c *Ctx) Execute(fn func(...string), commands ...func(*cobra.Command)) {
	for i := range commands {
		commands[i](c.cmd)
	}
	c.cmd.Use = c.name
	c.cmd.Version = c.version
	c.cmd.Run = func(cmd *cobra.Command, args []string) {
		for i := range c.conf {
			c.log(c.conf[i])
		}
		fn(args...)
	}
	// TODO implement app deploy config generator
	// for name, dpl := range c.deployers {
	// 	c.AddCommand(name, func(...string) {
	// 		if setter, ok := dpl.(types.DefaultSetter); ok {
	// 			setter.SetDefault()
	// 		}
	// 		filename := path.Join(c.root, name)
	// 		if err := dpl.Write(filename); err != nil {
	// 			panic(fmt.Errorf("init %s error: %v", name, err))
	// 		}
	// 	}, func(cmd *cobra.Command) {
	// 		cmd.Short = "init configuration for " + name
	// 	})
	// }
	if err := c.cmd.Execute(); err != nil {
		panic(err)
	}
}

func (c *Ctx) scan(rv reflect.Value) error {
	vars := env.NewVars(c.group(rv))

	if err := env.NewDecoder(vars).Decode(rv); err != nil {
		return err
	}
	c.vars = append(c.vars, vars)
	if _, err := env.NewEncoder(vars).Encode(rv); err != nil {
		return err
	}
	return nil
}

func (c *Ctx) marshal(rv reflect.Value) error {
	vars := env.LoadVarsFromEnviron(c.group(rv), os.Environ())
	if err := env.NewDecoder(vars).Decode(rv); err != nil {
		return err
	}
	return nil
}

func (c *Ctx) MarshalDefault() error {
	// TODO: add comment for each single config element
	m := map[string]string{
		consts.GoRuntimeEnv: consts.DevelopEnv,
	}
	for _, vars := range c.vars {
		for _, v := range vars.Values {
			if !v.Optional {
				m[v.Key(vars.Prefix)] = v.Value
			}
		}
	}

	return WriteYamlFile(path.Join(c.root, "./config/config.yml.template"), m)
}

func (c *Ctx) log(rv reflect.Value) {
	vars := env.NewVars(c.group(rv))
	if _, err := env.NewEncoder(vars).Encode(rv); err != nil {
		panic(err)
	}
	fmt.Printf("%s", string(vars.MaskBytes()))
}

type Marshaller func(v interface{}) ([]byte, error)

// group returns config group name
func (c *Ctx) group(rv reflect.Value) string {
	group := rv.Elem().Type().Name()
	if rv.Elem().Type().Implements(types.RTypeNamed) {
		group = rv.Elem().Interface().(types.Named).Name()
	}
	if group == "" {
		return strings.ToUpper(strings.Replace(c.name, "-", "_", -1))
	}
	return strings.ToUpper(strings.Replace(c.name+"__"+group, "-", "_", -1))
}
