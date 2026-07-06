package main

import (
	"fmt"
	"os"

	"github.com/SoFarCalm/blobgator/internal/config"
)

func main() {
	var s state
	cfg := config.ReadConfig()
	s.configPtr = &cfg

	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments provided")
		os.Exit(1)
	}

	commandName := args[1]
	commandArgs := args[1:]

	userCmd := command{
		name: commandName,
		args: commandArgs,
	}

	cmds.run(&s, userCmd)
	os.Exit(0)
}
