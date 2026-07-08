package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/SoFarCalm/blobgator/internal/config"
	"github.com/SoFarCalm/blobgator/internal/database"
	"github.com/SoFarCalm/blobgator/internal/rssfeed"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Error reading config %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		db:        dbQueries,
		configPtr: &cfg,
	}

	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmds.run(programState, command{name: cmdName, args: cmdArgs})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := rssfeed.Client{}
	client.FetchFeed(ctx, rssfeed.BaseURL)
	os.Exit(0)
}
