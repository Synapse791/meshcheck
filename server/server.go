package server

import (
	"fmt"
	"net/http"
	"github.com/Synapse791/meshcheck/config"
	"github.com/Synapse791/meshcheck/logger"
	"encoding/json"
)

type Server struct {
	Config config.AppConfig
}

type ServerResponse struct {
	Success             bool                `json:"success"`
	FailedConnections   []config.Connection	`json:"failed_connections"`
}

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

func (s Server) PingClients() *ServerResponse {

	return &ServerResponse{
		true,
		nil,
	}

}