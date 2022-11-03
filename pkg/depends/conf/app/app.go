package app

import (
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"

	"github.com/machinefi/w3bstream/pkg/depends/conf/deploy"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
)

type OptSetter = func(conf *Ctx)

func WithName(name string) OptSetter { return func(c *Ctx) { c.name = name } }

func WithVersion(version string) OptSetter { return func(c *Ctx) { c.version = version } }

func WithFeat(feat string) OptSetter { return func(c *Ctx) { c.feat = feat } }

func WithRoot(root string) OptSetter {
	_, filename, _, _ := runtime.Caller(1)
	return func(c *Ctx) {
		info, err := os.Stat(filepath.Join(".", "config"))
		if err == nil && info.IsDir() {
			c.root = "."
			return
		}
		c.root = filepath.Join(filepath.Dir(filename), root)
	}
}

func WithLogger(l log.Logger) OptSetter { return func(c *Ctx) { c.ctx = log.WithLogger(c.ctx, l) } }

func WithDeployer(deployers ...deploy.Deployer) OptSetter {
	return func(c *Ctx) {
		if c.deployers == nil {
			c.deployers = make(map[string]deploy.Deployer)
		}
		for _, dpl := range deployers {
			c.deployers[dpl.Name()] = dpl
		}
	}
}

func WriteYamlFile(filename string, v interface{}) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	root := filepath.Dir(filename)
	if root != "" {
		if err = os.MkdirAll(root, os.ModePerm); err != nil {
			return err
		}
	}
	return os.WriteFile(filename, data, os.ModePerm)
}
