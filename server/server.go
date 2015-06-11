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
	Success     bool                `json:"success"`
	Connections ResponseConnections `json:"connections"`
	Errors      []string            `json:"errors"`
}

type ResponseConnections struct {
	Successful []config.Connection `json:"successful"`
	Failed     []config.Connection `json:"failed"`
}

func GetInitMessage() string {
	return "Running in Server mode"
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SetConfig(configDir string) error {

	config, err := config.GetServerConfig(configDir)

	if err != nil {
		return err
	}

	s.Config = config

	return nil
}

func (s Server) Listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){

		logger.Info(fmt.Sprintf("Recieved call from %s", r.RemoteAddr))

		response := s.PingClients()

		output, err := json.Marshal(response)
		if err != nil {
			logger.Warning("Failed to encode response")
			return
		}

		fmt.Fprint(w, string(output))

	})

	logger.Info(fmt.Sprintf("Listening on port %s", s.Config.Port))

	http.ListenAndServe(s.Config.Port, nil)
}

func (s Server) PingClients() ServerResponse {

	var resp ServerResponse
	resp.Success = true

	for _, conn := range s.Config.Connections {

		address := fmt.Sprintf("http://%s:%s", conn.ToAddress, conn.Port)

		logger.Info(fmt.Sprintf("Calling %s", address))

		data, err := http.Get(address)
		if err != nil {
			msg := fmt.Sprintf("Failed to connect to %s", address)
			logger.Warning(msg)
			resp.Errors = append(resp.Errors, msg)
			resp.Success = false
			return resp
		}

		defer data.Body.Close()
		body, bodyErr := ioutil.ReadAll(data.Body)
		if bodyErr != nil {
			msg := fmt.Sprintf("Failed to read response body from %s", address)
			logger.Warning(msg)
			resp.Errors = append(resp.Errors, msg)
			resp.Success = false
			return resp
		}

		var cResp client.ClientResponse

		jsonErr := json.Unmarshal(body, &cResp)
		if jsonErr != nil {
			msg := fmt.Sprintf("Failed to decode client response from %s", address)
			logger.Warning(msg)
			resp.Errors = append(resp.Errors, msg)
			resp.Success = false
			return resp
		}

		s.ParseClientResponse(conn.ToAddress, cResp, &resp)

		logger.Info(fmt.Sprintf("%s responded: %t", address, cResp.Success))
	}

	return resp

}

func (s Server) ParseClientResponse(cAddr string, cResp client.ClientResponse, sResp *ServerResponse) {

	if cResp.Success == false {
		sResp.Success = false
	}

	for _, conn := range cResp.Connections {
		conn.FromAddress = cAddr
		if conn.Success == false {
			sResp.Connections.Failed = append(sResp.Connections.Failed, conn)
		} else {
			sResp.Connections.Successful = append(sResp.Connections.Successful, conn)
		}
	}

}