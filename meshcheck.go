package main

import (
	"os"
	"fmt"
	"github.com/Synapse791/meshcheck/client"
	"github.com/Synapse791/meshcheck/server"
)

func main() {
	args := os.Args[1:]

	mode := args[0]

	if mode == "client" {
		fmt.Println(client.GetHelp())
	} else if mode == "server" {
		fmt.Println(server.GetHelp())
	} else {
		fmt.Println("Invalid mode [client|server]")
	}
}