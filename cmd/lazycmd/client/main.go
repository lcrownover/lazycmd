package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lcrownover/lazycmd/internal/lazycmd"
)

func main() {

	fTarget := flag.String("target", "", "target host")
	fCommand := flag.String("command", "", "command to run")
	flag.Parse()

	if *fTarget == "" || *fCommand == "" {
		fmt.Println("You must provide both a -target and -command")
		os.Exit(1)
	}

	target := lazycmd.CleanseTarget(*fTarget)
	command := lazycmd.CleanseCommand(*fCommand)

	p := lazycmd.NewProducer(target)
	defer p.Close()

	p.SendMessage(target, command)
}
