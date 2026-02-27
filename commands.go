package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.name]

	if !ok {
		return fmt.Errorf("command %s does not exist!\n", cmd.name)
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if _, ok := c.registeredCommands[name]; ok {
		return fmt.Errorf("command already exists!\n")
	}

	c.registeredCommands[name] = f

	return nil
}
