package server

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/Synapse791/meshcheck/config"
	"github.com/Synapse791/meshcheck/logger"
	"encoding/json"
	"github.com/Synapse791/meshcheck/client"
)

type Server struct {
	Config config.AppConfig
}

type ServerResponse struct {
	Success             bool                      `json:"success"`
	FailedConnections   []config.FailedConnection `json:"failed_connections"`
	Errors              []string                  `json:"errors"`
}

//TODO Implement the following response structure
//type ServerResponse struct {
//	Success     bool                `json:"success"`
//	Connections ResponseConnections `json:"connections"`
//	Errors      []string            `json:"errors"`
//}
//
//type ResponseConnections struct {
//	Successful []config.Connection `json:"successful"`
//	Failed     []config.Connection `json:"failed"`
//}

func GetInitMessage() string {
	return "Server mode set"
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SetConfig(configDir string) bool {

	config, err := config.GetServerConfig(configDir)

	if err != nil {
		return false
	}

	s.Config = config

	return true
}

func (s Server) Listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){

		response := s.PingClients()

		if response.Success == false {
			//TODO set Status Code
		}

		output, err := json.Marshal(response)
		if err != nil {
			logger.Warning("Failed to encode response")
			return
		}

		fmt.Fprint(w, string(output))

	})

	logger.Info("Listening on port " + s.Config.Port)

	http.ListenAndServe(s.Config.Port, nil)
}

func (s Server) PingClients() ServerResponse {

	var resp ServerResponse
	resp.Success = true

	for _, conn := range s.Config.Connections {

		address := "http://" + conn.IpAddress + ":" + conn.Port

		logger.Info("Calling " + address)

		data, err := http.Get(address)
		if err != nil {
			msg := "Failed to connect to " + address
			logger.Warning(msg)
			resp.Errors = append(resp.Errors, msg)
			resp.Success = false
			return resp
		}

		logger.Info("Got response")

		defer data.Body.Close()
		body, bodyErr := ioutil.ReadAll(data.Body)
		if bodyErr != nil {
			msg := "Failed to read response body from " + address
			logger.Warning(msg)
			resp.Errors = append(resp.Errors, msg)
			resp.Success = false
			return resp
		}

		var cResp client.ClientResponse

		jsonErr := json.Unmarshal(body, &cResp)
		if jsonErr != nil {
			msg := "Failed to decode client response from " + address
			logger.Warning(msg)
			resp.Errors = append(resp.Errors, msg)
			resp.Success = false
			return resp
		}

		s.ParseClientResponse(address, cResp, &resp)

		logger.Info(address + " responded: \n" + string(body))
	}

	return resp

}

func (s Server) ParseClientResponse(cAddr string, cResp client.ClientResponse, sResp *ServerResponse) {

	if cResp.Success == false {
		sResp.Success = false
		for _, conn := range cResp.Connections {
			if conn.Success == false {

				var failure config.FailedConnection

				failure.ClientAddress     = cAddr
				failure.ConnectionAddress = conn.IpAddress
				failure.Port              = conn.Port

				sResp.FailedConnections = append(sResp.FailedConnections, failure)
			}
		}
	}

}