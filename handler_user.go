package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SoFarCalm/blobgator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		fmt.Println("Please provide a username")
		os.Exit(1)
	}

	name := cmd.args[0]

	returnedUser, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		log.Fatalf("user %s does not exist", name)
		os.Exit(1)
	}

	if len(returnedUser) == 0 {
		fmt.Printf("User %s dosen't exists\n", name)
		os.Exit(1)
	}

	err = s.configPtr.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Println("User has been set")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		fmt.Println("Please provide a username")
		os.Exit(1)
	}

	name := cmd.args[0]

	returnedUser, _ := s.db.GetUser(context.Background(), name)

	if returnedUser == name {
		fmt.Printf("User %s already exists\n", name)
		os.Exit(1)
	}

	c := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	createdUser, err := s.db.CreateUser(context.Background(), c)
	if err != nil {
		log.Fatalf("Error when creating user %v\n", err)
	}

	err = s.configPtr.SetUser(name)
	if err != nil {
		log.Fatalf("Error when setting user %v\n", err)
	}

	fmt.Println("User has been created as follows")
	fmt.Printf("- %s\n", createdUser.ID)
	fmt.Printf("- %s\n", createdUser.CreatedAt)
	fmt.Printf("- %s\n", createdUser.UpdatedAt)
	fmt.Printf("- %s\n", createdUser.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		log.Fatalf("Error when deleting from users table: %v\n", err)
	}

	fmt.Println("All users have been deleted")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("Error when fecthing uesrs %v\n", err)
	}

	for i := 0; i < len(users); i++ {
		if users[i].Name == s.configPtr.CurrentUsername {
			fmt.Printf("- %s (current)\n", users[i].Name)
		} else {
			fmt.Printf("- %s\n", users[i].Name)
		}
	}

	return nil
}
