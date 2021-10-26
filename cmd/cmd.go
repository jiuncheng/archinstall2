package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
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
	// args := strings.Split(c.args, " ")
	// args := strings.Fields(c.args)

	r := regexp.MustCompile(`[^\s"]+|"([^"]*)"`)
	args := r.FindAllString(c.args, -1)

	var trimmedArgs []string
	for _, arg := range args {
		trimmedArgs = append(trimmedArgs, strings.Trim(arg, "\""))
	}

	fmt.Printf("\n%s\n", c.desc)
	fmt.Printf("Executing command : %s\n", c.args)

	for _, arg := range args {
		log.Println(arg)
	}

	var cmd *exec.Cmd
	if len(args) < 2 {
		cmd = exec.Command(trimmedArgs[0])
	} else {
		cmd = exec.Command(trimmedArgs[0], trimmedArgs[1:]...)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	_, errStr := stdoutBuf.String(), stderrBuf.String()
	fmt.Printf("%s", errStr)
	log.Println("Operation done.")
	return err

	// output, err := cmd.CombinedOutput()
	// fmt.Println(string(output))
	// log.Println("Operation done.")
	// return err
}

func (c *Command) Output() ([]byte, error) {
	args := strings.Split(c.args, " ")

	fmt.Printf("\n%s\n", c.desc)
	fmt.Printf("\nExecuting command : %s\n", c.args)

	var cmd *exec.Cmd
	if len(args) < 2 {
		cmd = exec.Command(args[0])
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}
	return cmd.Output()
}
