package main

import (
	"log"
	"os"
	"database/sql"

	"github.com/gclinoz/Gator-go/internal/config"
	"github.com/gclinoz/Gator-go/internal/database"

	_ "github.com/lib/pq"
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
	cmds.register("users", handlerListUser)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeed)
	cmds.register("follow", middlewareLoggedIn(handlerAddFollow))
	cmds.register("following", middlewareLoggedIn(handlerListFollow))
	cmds.register("unfollow", middlewareLoggedIn(handlerDelFollow))
	cmds.register("browse", middlewareLoggedIn(handlerPost))

	if len(os.Args) < 2 {
		log.Fatal("no command provided")
	}

	cmdsInput := command{name: os.Args[1], args: os.Args[2:]}
	err = cmds.run(&st, cmdsInput)
	if err != nil {
		log.Fatal(err)
	}
}
