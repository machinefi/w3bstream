package cmd

import (
	"bytes"
	"fmt"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func NewCommand(name string, args string) *Command {
	cmd := &Command{Name: name}
	if _args := strings.Split(args, " "); len(args) > 0 {
		cmd.Args = _args
	}
	return cmd
}

func (c Command) MarshalText() ([]byte, error) {
	return []byte(strings.Join(append([]string{c.Name}, c.Args...), " ")), nil
}

func (c *Command) UnmarshalText(data []byte) error {
	parts := bytes.Split(data, []byte(" "))
	if len(parts) == 0 {
		return fmt.Errorf("unexpected command line: %s", string(data))
	}
	c.Name = string(parts[0])
	c.Args = c.Args[0:0]

	for _, v := range parts[1:] {
		c.Args = append(c.Args, string(v))
	}
	return nil
}
