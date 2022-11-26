package main

import (
	"fmt"

	"github.com/lcrownover/lazycmd/internal/lazycmd"
)

func main() {
	h := lazycmd.GetHostname()
	myHostname := lazycmd.CleanseTarget(h)

	c := lazycmd.NewConsumer(myHostname)
	defer c.Close()

	for {
		msg, some := c.GetMessage()
		if some {
			fmt.Printf("%s\n", msg)
			cmd := lazycmd.NewCommand(msg)
			out, err := cmd.Execute()
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
			fmt.Printf("output: %s\n", out)
		}
	}
}
