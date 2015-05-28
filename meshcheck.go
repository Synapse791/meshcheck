package main

import (
	"os"
	"github.com/Synapse791/meshcheck/config"
	"github.com/Synapse791/meshcheck/logger"
	"github.com/Synapse791/meshcheck/client"
	"github.com/Synapse791/meshcheck/server"
)

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		logger.Fatal("Missing arguments. Must pass mode and config directory. Example: meshcheck client /etc/meshcheck/conf")
	}

	mode := args[0]
	configDir := args[1]

	if mode == "client" {
		logger.Info(client.GetInitMessage())

		config, err := config.GetClientConfig(configDir)

		if err != nil {
			logger.Fatal("Failed to read config files")
		} else {
			logger.Info("Config set")
		}

		c := client.NewClient()

		c.Config = config
		c.Listen()

	} else if mode == "server" {
		logger.Info(server.GetHelp())
	} else {
		logger.Fatal("Invalid mode {client|server}")
	}

}