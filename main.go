package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Brime/gatorcli/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	s := state{
		cfg: &conf,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Errorf("invalid command")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	cmds.run(s, cmd)

	conf, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", conf)
}
