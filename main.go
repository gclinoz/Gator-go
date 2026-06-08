package main

import _"github.com/lib/pq"

import (
	"log"
	"os"
	"database/sql"

	"github.com/gclinoz/Gator-go/internal/config"
	"github.com/gclinoz/Gator-go/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	defer db.Close()

	dbQueries := database.New(db)
	st := state{
		cfg:	&cfg,
		db:		dbQueries,
		
	}
	cmds := commands{
		utils: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegis)
	cmds.register("reset", handlerReset)

	if len(os.Args) < 2 {
		log.Fatal("no command provided")
	}

	cmdsInput := command{name: os.Args[1], args: os.Args[2:]}
	err = cmds.run(&st, cmdsInput)
	if err != nil {
		log.Fatal(err)
	}
}
