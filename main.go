package main

import (
	"log"
	"os"

	"github.com/roshankaranth/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {

	cfg, err := config.Read()

	if err != nil {
		log.Fatalf("%v", err)
	}

	programState := &state{
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	err = cmds.register("login", handlerLogin)

	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("Not enough args provided!\n")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
