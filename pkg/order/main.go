package main

import (
	"os"
	"os/signal"
	"pkg/pkg/order/cmd/server"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(0)
}

func main() {
	startTime := time.Now()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.WithField("uptime", time.Since(startTime).String()).
			WithField("signal", sig.String()).
			Error("interrupt signal detected")
		os.Exit(0)
	}()

	s := server.Server{
		Port: 8090,
	}

	log.Fatal(s.Serve())
}