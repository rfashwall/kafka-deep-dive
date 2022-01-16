package server

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/rfashwall/kafka-deep-dive/pkg/post/internal/handlers"
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
