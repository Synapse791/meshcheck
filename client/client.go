package client

import (
	"fmt"
	"net"
	"net/http"
	"github.com/Synapse791/meshcheck/config"
	"github.com/Synapse791/meshcheck/logger"
	"encoding/json"
)

type Client struct {
	Config  config.AppConfig
}

type ClientResponse struct {
	Success         bool                `json:"success"`
	Connections     []config.Connection `json:"connections"`
}

func GetInitMessage() string {
	return "Client mode set"
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SetConfig(configDir string) error {

	config, err := config.GetClientConfig(configDir)

	if err != nil {
		return err
	}

	c.Config = config

	return nil

}

func (c Client) Listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){

		logger.Info(fmt.Sprintf("Recieved call from %s", r.RemoteAddr))

		response := c.ScanPorts()

		output, err := json.Marshal(response)
		if err != nil {
			logger.Warning("Failed to encode response")
			return
		}

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprint(w, string(output))

	})

	logger.Info(fmt.Sprintf("Listening on port %s", c.Config.Port))

	http.ListenAndServe(c.Config.Port, nil)
}

func (c *Client) ScanPorts() *ClientResponse {

	resp := &ClientResponse {
		true,
		c.Config.Connections,
	}

	logger.Info("Running connection list")

	for count, connection := range c.Config.Connections {

		address := fmt.Sprintf("%s:%s", connection.ToAddress, connection.Port)

		check := c.CheckConnection(address)

		if check {
			resp.Connections[count].Success = true
			logger.Info(fmt.Sprintf("%s -> OK", address))
		} else {
			logger.Warning(fmt.Sprintf("%s -> FAILED", address))
		}

	}

	for _, check := range resp.Connections {
		if check.Success == false {
			resp.Success = false
		}
	}

	return resp

}

func (c *Client) CheckConnection(addr string) bool {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return false
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return false
	}

	defer conn.Close()

	return true
}
