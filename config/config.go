package config

import (
	"os"
	"io/ioutil"
	"bufio"
	"strings"
	"regexp"
	"github.com/Synapse791/meshcheck/logger"
	"fmt"
)

type AppConfig struct {
	Port            string
	Mode            string
	Connections     []Connection
}

type Connection struct {
	Success     bool    `json:"success"`
	FromAddress string  `json:"from_address"`
	ToAddress   string  `json:"to_address"`
	Port        string  `json:"port"`
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

	filePaths["connections"] = fmt.Sprintf("%sconnections", dir)
	filePaths["port"]        = fmt.Sprintf("%sport", dir)

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
		return fmt.Errorf("os: failed to stat file %s", file)
	}

	connFile, connErr := os.Open(file)
	if connErr != nil {
		return fmt.Errorf("os: failed to open file %s", file)
	}

	defer connFile.Close()

	scanner := bufio.NewScanner(connFile)

	for scanner.Scan() {
		line := scanner.Text()

		match, regErr := regexp.MatchString("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}:[0-9]+", line)
		if regErr != nil {
			return fmt.Errorf("regexp: error processing regex")
		}

		if !match {
			logger.Fatal(fmt.Sprintf("Invalid connection (%s). Connections should be a list of IP:PORT combinations", line))
		}

		var conn Connection

		data := strings.Split(line, ":")

		conn.ToAddress = data[0]
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

		config.Port = fmt.Sprintf(":%s", port)
		config.Port = strings.TrimSpace(config.Port)

	}

	return nil
}