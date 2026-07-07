package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/SoFarCalm/blobgator/internal/config"
	"github.com/SoFarCalm/blobgator/internal/database"
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

	// userCmd := command{
	// 	name: cmdName,
	// 	args: cmdArgs,
	// }

	cmds.run(programState, command{name: cmdName, args: cmdArgs})
	os.Exit(0)
}
