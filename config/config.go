package config

import (
	"os"
	"io/ioutil"
	"bufio"
	"strings"
	"regexp"
	"github.com/Synapse791/meshcheck/logger"
)

type AppConfig struct {
	Port            string
	Mode            string
	Connections     []Connection
}

type Connection struct {
	Success     bool    `json:"success"`
	IpAddress   string  `json:"ip_address"`
	Port        string  `json:"port"`
}

type FailedConnection struct {
	ClientAddress     string  `json:"client_address"`
	ConnectionAddress string  `json:"connection_address"`
	Port              string  `json:"port"`
}

func GetClientConfig(dir string) (AppConfig, error) {

	var config AppConfig

	config.Mode = "client"

	if err := SetConfig(dir, &config); err != nil {
		return config, err
	}

	return config, nil

}

func GetServerConfig(dir string) (AppConfig, error) {

	var config AppConfig

	config.Mode = "server"

	if err := SetConfig(dir, &config); err != nil {
		return config, err
	}

	return config, nil

}

func SetConfig(dir string, config *AppConfig) error {
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

	filePaths := make(map[string]string)

	filePaths["connections"]	= dir + "connections"
	filePaths["port"] 			= dir + "port"

	if err := ReadConnectionConfig(filePaths["connections"], config); err != nil {
		return err
	}

	if err := ReadPortConfig(filePaths["port"], config); err != nil {
		return err
	}

	return nil
}

func ReadConnectionConfig(file string, config *AppConfig) error {
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

		var conn Connection

		data := strings.Split(line, ":")

		conn.IpAddress = data[0]
		conn.Port      = data[1]

		config.Connections = append(config.Connections, conn)

	}

	return nil
}

func ReadPortConfig(file string, config *AppConfig) error {

	if _, err := os.Stat(file); os.IsNotExist(err) {

		if config.Mode == "client" {
			logger.Info("Setting default client port (6600)")
			config.Port = ":6600"
		} else if config.Mode == "server" {
			logger.Info("Setting default server port (6800)")
			config.Port = ":6800"
		} else {
			logger.Fatal("Unkown Mode")
		}


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