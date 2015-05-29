package config

import (
	"os"
	"io/ioutil"
	"bufio"
	"strings"
	"regexp"
	"github.com/Synapse791/meshcheck/logger"
)

type ClientConfig struct {
	Port			string
	Connections		[]string
}

func GetClientConfig(dir string) (ClientConfig, error) {

	var config ClientConfig

	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

	filePaths := make(map[string]string)

	filePaths["connections"]	= dir + "connections"
	filePaths["port"] 			= dir + "port"

	if err := ReadClientConnectionConfig(filePaths["connections"], &config); err != nil {
		return config, err
	}

	if err := ReadClientPortConfig(filePaths["port"], &config); err != nil {
		return config, err
	}

	return config, nil

}

func ReadClientConnectionConfig(file string, config *ClientConfig) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return err
	}

	connFile, connErr := os.Open(file)
	if connErr != nil {
		return connErr
	}

	defer connFile.Close()

	scanner := bufio.NewScanner(connFile)

	for scanner.Scan() {
		line := scanner.Text()

		match, _ := regexp.MatchString("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}:[0-9]+", line)

		if !match {
			logger.Fatal("Invalid connection ("+line+"). Connections should be a list of IP:PORT combinations")
		}

		config.Connections = append(config.Connections, line)

	}

	return nil
}

func ReadClientPortConfig(file string, config *ClientConfig) error {

	if _, err := os.Stat(file); os.IsNotExist(err) {

		logger.Info("Setting default port (6600)")
		config.Port = ":6600"

	} else {

		port, portErr := ioutil.ReadFile(file)
		if portErr != nil {
			return portErr
		}

		config.Port = ":" + string(port)
		config.Port = strings.TrimSpace(config.Port)

	}

	return nil
}