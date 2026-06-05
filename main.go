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
		ptr: &cfg,
	}
	cmds := commands{
		utils: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatalln("no command provided")
	}

	argBox := []string{}
	for idx, val := range os.Args {
		if idx > 1 {
			argBox = append(argBox, val)
		}
	}

	cmdsInput := command{
		name: os.Args[1],
		args: argBox,
	}
	err = cmds.run(&st, cmdsInput)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}
}
