package client

import (
	"os"
	"bufio"
	"strings"
)

type Config struct {
	FilePath		string
	Connections		[]Connection
}

type Connection struct {
	IpAddress	string
	Port		string
}

func ReadConfigFile(dir string) (Config, error) {

	var config Config
	config.FilePath = dir + "/connections"

	if _, err := os.Stat(config.FilePath); os.IsNotExist(err) {
		return config, err
	}

	file, fileErr := os.Open(config.FilePath)

	if fileErr != nil {
		return config, fileErr
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var conn Connection

		data := strings.Split(scanner.Text(), ":")

		conn.IpAddress = data[0]
		conn.Port      = data[1]

		config.Connections = append(config.Connections, conn)

	}

	return config, nil
}