package main

import "fmt"

type Command struct {
	name string
	args []string
}

type Commands struct {
	cmd map[string]func(*State, Command) error
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.cmd[name] = f
}

func (c *Commands) run(s *State, cmd Command) error {
	exec, exist := c.cmd[cmd.name]
	if !exist {
		return fmt.Errorf("command %v don't exist", cmd.name)
	}
	return exec(s, cmd)
}
