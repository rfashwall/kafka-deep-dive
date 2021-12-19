package server

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"pkg/pkg/order/pkg/handlers"
)

type Server struct {
	Port int
}

func (s *Server) Serve() error {
	http.HandleFunc("/orders/health", handlers.Health)
	address := fmt.Sprintf(":%d", s.Port)

	log.WithField("address", address).Info("server starting")

	return http.ListenAndServe(address, nil)
}
