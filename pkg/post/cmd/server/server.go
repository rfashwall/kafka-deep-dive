package server

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"pkg/pkg/post/pkg/handlers"
)

type Server struct {
	Port int
}

func (s *Server) Serve() error {
	http.HandleFunc("/posts/health", handlers.Health)
	address := fmt.Sprintf(":%d", s.Port)

	log.WithField("address", address).Info("server starting")

	return http.ListenAndServe(address, nil)
}