package section_config_test

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/base/cmd"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	. "github.com/machinefi/w3bstream/pkg/depends/conf/section_config"
)

type Config struct {
	Section
	Command      cmd.Command             `name:"command"`
	AutoStart    bool                    `name:"autostart"`
	AutoRestart  bool                    `name:"autorestart"`
	StartSecs    types.Second            `name:"startsecs"`
	StartRetries int                     `name:"startretries"`
	User         string                  `name:"user"`
	Priority     int                     `name:"priority"`
	ExitCodes    types.CommaSplitInts    `name:"exitcodes"`
	StopSignal   types.Signal            `name:"stopsignal"`
	Envs         types.CommaSplitStrings `name:"environment"`
	*LogConfig
}

func (c *Config) GetSection() *Section { return &c.Section }

func (c *Config) SetSection(k, v string) { c.Name, c.Value = k, v }

type LogConfig struct {
	StdoutLogFile         string   `name:"stdout_logfile"`
	StdoutLogFileMaxBytes types.MB `name:"stdout_logfile_maxbytes"`
	StdoutLogFileBackups  int      `name:"stdout_logfile_backups"`
}

func TestLoadFile(t *testing.T) {
	var (
		cwd, _ = os.Getwd()
		path   = filepath.Join(cwd, "__example__", "demo.conf")
		cfg    = &Config{
			Section:      Section{Name: "program", Value: "crond"},
			Command:      cmd.Command{Name: "/opt/iTR/core/sbin/crond", Args: []string{"-n"}},
			AutoStart:    true,
			AutoRestart:  true,
			StartSecs:    types.Second(3),
			StartRetries: 10,
			User:         "root",
			Priority:     200,
			ExitCodes:    types.CommaSplitInts{0},
			StopSignal:   types.Signal(syscall.SIGTERM),
			Envs:         types.CommaSplitStrings{`PATH="/opt/iTR/core/bin:/opt/iTR/core/sbin"`},
			LogConfig: &LogConfig{
				StdoutLogFile:         "/opt/iTR/core/var/log/supervisor/crond.log",
				StdoutLogFileMaxBytes: types.MB(1),
				StdoutLogFileBackups:  3,
			},
		}
		cpy = &Config{}
	)

	if err := NewEncoder('=').MarshalToFile(cfg, path); err != nil {
		t.Log(err)
		return
	}

	if err := NewDecoder('=').UnmarshalFromFile(cpy, path); err != nil {
		t.Log(err)
		return
	}

	NewWithT(t).Expect(cpy).To(Equal(cfg))

}
