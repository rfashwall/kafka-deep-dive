package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/rfashwall/kafka-deep-dive/pkg/order/internal/handlers"
)

type Server struct {
	Port int
}

func handler(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *Server) Serve() error {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/orders/health", handlers.Health).Methods(http.MethodGet)
	r.HandleFunc("/orders", handlers.Create).Methods(http.MethodPost)

	address := fmt.Sprintf(":%d", s.Port)

	log.WithField("address", address).Info("server starting")

	return http.ListenAndServe(address, r)
}
