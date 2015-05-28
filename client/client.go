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
	Config		config.ClientConfig
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

		response := c.ScanPorts()

		output, err := json.Marshal(response)
		if err != nil {
			logger.Warning("Failed to encode response")
			return
		}

		fmt.Fprint(w, string(output))

	})

	logger.Info("Listening on port " + c.Config.Port)

	http.ListenAndServe(c.Config.Port, nil)
}

func (c *Client) ScanPorts() *Response {

	resp := &Response{}

	for _, address := range c.Config.Connections {
		logger.Info("Checking " + address)
		check := c.CheckConnection(address)

		if !check {
			resp.Errors = append(resp.Errors, "Connection " + address + " failed")
		}
	}

	if len(resp.Errors) > 0 {
		resp.Success = false
	} else {
		resp.Success = true
	}

	return resp

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