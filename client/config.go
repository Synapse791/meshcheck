package client

import (
	"os"
	"encoding/json"
)

type Config struct {
	FilePath		string
	Connections		[]Connection
}

type Connection struct {
	IpAddress	string	`json:"ip"`
	Port		int		`json:"port"`
}

func ReadConfigFile() (Config, error) {

	var config Config
	config.FilePath = "/etc/meshcheck/connections.json"

	if _, err := os.Stat(config.FilePath); os.IsNotExist(err) {
		return config, err
	}

	// TODO: error handling
	file, _ := os.Open(config.FilePath)

	decoder := json.NewDecoder(file)

	decoder.Decode(&config)

	return config, nil
}