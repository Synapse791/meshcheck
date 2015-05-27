package client

import (
	"fmt"
	"net"
	"net/http"
	"github.com/Synapse791/meshcheck/logger"
	"encoding/json"
)

type Client struct {
	Config		Config
	Response	Response
}

type Response struct {
	Success	bool		`json:"success"`
	Errors  []string	`json:"errors"`
}

func GetInitMessage() string {
	return "Client mode set"
}

func NewClient() *Client {
	return &Client{}
}

func (c Client) Listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){


		c.ScanPorts()

		output, err := json.Marshal(c.Response)
		if err != nil {
			logger.Warning("Failed to encode response")
			return
		}

		fmt.Fprint(w, string(output))

		c.Response.Errors = []string{}
	})

	logger.Info("Listening on port 6600")

	http.ListenAndServe(":6600", nil)
}

func (c *Client) ScanPorts() {

	for _, connection := range c.Config.Connections {
		address := c.BuildAddress(connection.IpAddress, connection.Port)

		logger.Info("Checking " + address)
		check := c.CheckConnection(address)

		if !check {
			c.Response.Errors = append(c.Response.Errors, "Connection " + connection.IpAddress + ":" + connection.Port + " failed")
		}
	}

	if len(c.Response.Errors) > 0 {
		c.Response.Success = false
	} else {
		c.Response.Success = true
	}

}

func (c *Client) CheckConnection(addr string) bool {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		logger.Warning("Failed to resolve " + addr)
		return false
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logger.Warning("Failed to connect to " + addr)
		return false
	}

	defer conn.Close()

	return true
}

func (c *Client) BuildAddress(ip string, port string) string {
	return ip + ":" + port
}