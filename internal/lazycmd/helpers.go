package lazycmd

import (
	"fmt"
	"os"
	"strings"
)

func GetHostname() string {
	h, err := os.Hostname()
	if err != nil {
		fmt.Println("failed to get hostname")
		panic(err)
	}
	return h
}

func CleanseTarget(t string) string {
	// This should run on both the server and client.
	// That way, the hostnames will always match.
	ch := strings.ToLower(t)
	return ch
}

func CleanseCommand(c string) string {
	// Cleaning the command if needed
	return c
}
