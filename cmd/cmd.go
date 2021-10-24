package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

type Command struct {
	args string
	desc string
}

func NewCmd(args string) *Command {
	return &Command{
		args: args,
	}
}

func (c *Command) GetArgs() string {
	return c.args
}

func (c *Command) GetDesc() string {
	return c.desc
}

func (c *Command) SetArgs(args string) *Command {
	c.args = args
	return c
}

func (c *Command) SetDesc(desc string) *Command {
	c.desc = desc
	return c
}

func (c *Command) Run() error {
	args := strings.Split(c.args, " ")

	fmt.Printf("\n%s\n", c.desc)
	fmt.Printf("\nExecuting command : %s\n", c.args)

	cmd := exec.Command(args[0], args...)
	return cmd.Run()
}
