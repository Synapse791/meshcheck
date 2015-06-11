package main

import (
	"os"
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

		c := client.NewClient()

		if configErr := c.SetConfig(configDir); configErr != nil {
			logger.Fatal(configErr.Error())
		} else {
			logger.Info("Config set")
		}

		c.Listen()

	} else if mode == "server" {
		logger.Info(server.GetInitMessage())

		s := server.NewServer()

		if configErr := s.SetConfig(configDir); configErr != nil {
			logger.Fatal(configErr.Error())
		} else {
			logger.Info("Config set")
		}

		s.Listen()

	} else {
		logger.Fatal("Invalid mode {client|server}")
	}

}