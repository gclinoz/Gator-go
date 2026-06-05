package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gclinoz/Gator-go/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %v\n", cfg)

	st := state{
		cfg: &cfg,
	}
	cmds := commands{
		utils: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("no command provided")
	}

	cmdsInput := command{name: os.Args[1], args: os.Args[2:]}
	err = cmds.run(&st, cmdsInput)
	if err != nil {
		log.Fatal(err)
	}
}
