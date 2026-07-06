package main

import (
	"fmt"
	"os"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) <= 1 {
		fmt.Println("Please provide a username")
		os.Exit(1)
	}

	name := cmd.args[1]
	err := s.configPtr.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Println("User has been set")

	return nil
}
