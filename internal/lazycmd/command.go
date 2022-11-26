package lazycmd

import (
	"fmt"
	"os/exec"
)

type Command struct {
	Command string
}

func NewCommand(command string) *Command {
	return &Command{
		Command: command,
	}
}

func (c *Command) Execute() (string, error) {
	fmt.Printf("executing command: %s\n", c.Command)
	out, err := exec.Command(c.Command).CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
