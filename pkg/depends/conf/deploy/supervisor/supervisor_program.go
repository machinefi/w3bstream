package supervisor

import (
	"syscall"

	"github.com/machinefi/w3bstream/pkg/depends/base/cmd"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/section_config"
)

type Program struct {
	section_config.Section
	// Command which will be run when supervisord started
	Command cmd.Command `name:"command"`
	// AutoStart if true, program will start automatically when supervisor start
	AutoStart bool `name:"autostart"`
	// AutoRestart Specifies if supervisord should automatically restart program
	AutoRestart bool `name:"autorestart"`
	// StartSecs The total seconds which the program needs to stay running after
	// a startup to consider the start successful (moving the process from the
	// `STARTING` state to the `RUNNING` state)
	StartSecs types.Second `name:"startsecs"`
	// StartRetries The number of serial failure attempts before program started
	StartRetries int `name:"startretries"`
	// User Instruct supervisord to switch users to this UNIX user account
	// before doing any meaningful processing
	User string `name:"user"`
	// Priority the relative priority of the program in the start and shutdown ordering
	Priority int `name:"priority"`
	// ExitCodes The list of expected exit codes for this program used with `autorestart`.
	ExitCodes types.CommaSplitInts `name:"exitcodes"`
	//
	// StopSignal The signal used to kill the program. It is normally one of:
	// `TERM`, `HUP`, `INT`, `QUIT`, `KILL`, `USR1`, or `USR2`
	StopSignal types.Signal `name:"stopsignal"`
	// Envs program's env vars
	Envs types.CommaSplitStrings `name:"environment"`
	// StdoutLogFile Put process stdout output in this file
	StdoutLogFile string `name:"stdout_logfile"`
	// StdoutLogFileMaxBytes The maximum number of bytes that may be consumed by
	// `StdoutLogFile`
	StdoutLogFileMaxBytes types.MB `name:"stdout_logfile_maxbytes"`
	// StdLogFileBackups The number of `stdout_logfile` backups to keep around
	// resulting from process stdout log file rotation
	StdoutLogFileBackups int `name:"stdout_logfile_backups"`
}

func (c Program) GetSection() *section_config.Section { return &c.Section }

func (c *Program) SetSection(name, val string) {
	c.Section.Name, c.Section.Value = name, val
}

func (c *Program) SetDefault() {
	if c.Priority == 0 {
		c.Priority = 500
	}
	if c.User == "" {
		c.User = "root"
	}
	if c.ExitCodes == nil {
		c.ExitCodes = types.CommaSplitInts{0}
	}
	if c.StopSignal == 0 {
		c.StopSignal = types.Signal(syscall.SIGTERM)
	}
}

func (c *Program) Write(filename string) error {
	return section_config.NewEncoder('=').MarshalToFile(c, filename)
}

func (c *Program) Name() string { return "supervisor.conf" }

func (c *Program) Bytes() ([]byte, error) {
	return section_config.NewEncoder('=').Marshal(c)
}
