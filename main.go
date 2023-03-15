package main

import (
	"github.com/seriallink/tsb/cmd"
	"log"
)

func main() {
	command := cmd.InitCmd()
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
