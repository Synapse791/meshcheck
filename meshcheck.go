package main

import (
	"os"
	"github.com/Synapse791/meshcheck/logger"
	"github.com/Synapse791/meshcheck/client"
	"github.com/Synapse791/meshcheck/server"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		logger.Fatal("Must specify client or server mode")
	}

	mode := args[0]

	if mode == "client" {
		logger.Info(client.GetInitMessage())
		TempFunc()
	} else if mode == "server" {
		logger.Info(server.GetHelp())
	} else {
		logger.Fatal("Invalid mode {client|server}")
	}
}

func TempFunc() {
	config, err := client.ReadConfigFile()

	if err != nil {
		logger.Fatal("Failed to read config file '" + config.FilePath + "'")
	} else {
		logger.Info("Config set")
	}

	c := client.NewClient()

	c.Config = config
	c.Listen()
}