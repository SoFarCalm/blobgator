package main

import (
	"github.com/SoFarCalm/blobgator/internal/config"
	"github.com/SoFarCalm/blobgator/internal/database"
)

type state struct {
	db        *database.Queries
	configPtr *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	err := c.commandMap[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}
