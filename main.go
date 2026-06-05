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
<<<<<<< HEAD
		ptr: &cfg,
=======
		cfg: &cfg,
>>>>>>> 5cafcae (CH1-L3: Commands)
	}
	cmds := commands{
		utils: make(map[string]func(*state, command) error),
	}
<<<<<<< HEAD

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
=======
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("no command provided")
	}

	cmdsInput := command{name: os.Args[1], args: os.Args[2:]}
	err = cmds.run(&st, cmdsInput)
	if err != nil {
		log.Fatal(err)
>>>>>>> 5cafcae (CH1-L3: Commands)
	}
}
