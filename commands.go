package main

import (
	"errors"
)

type command struct {
	name      string
	arguments []string
}

type commands struct {
	command map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.command == nil {
		c.command = make(map[string]func(*state, command) error)
	}
	c.command[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, exists := c.command[cmd.name]
	if !exists {
		return errors.New("command not found")
	}
	return f(s, cmd)
}
