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

	for _, fileCheck := range filePaths {
		if _, err := os.Stat(fileCheck); os.IsNotExist(err) {
			return config, err
		}
	}

	connFile, connErr := os.Open(filePaths["connections"])
	if connErr != nil {
		return config, connErr
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

	port, portErr := ioutil.ReadFile(filePaths["port"])
	if portErr != nil {
		return config, portErr
	}

	config.Port = ":" + string(port)
	config.Port = strings.TrimSpace(config.Port)

	return config, nil

}