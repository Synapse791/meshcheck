package main

import (
	"os"
	"github.com/Synapse791/meshcheck/client"
	"github.com/Synapse791/meshcheck/server"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		LogFatal("Must specify client or server mode")
	}

	mode := args[0]

	if mode == "client" {
		LogInfo(client.GetHelp())
		TempFunc()
	} else if mode == "server" {
		LogInfo(server.GetHelp())
	} else {
		LogFatal("Invalid mode {client|server}")
	}
}

func TempFunc() {
	config, err := client.ReadConfigFile()

	if err != nil {
		LogFatal("failed to read config file '" + config.FilePath + "'")
	}

	LogInfo(config.Connections[0].IpAddress)

}