package main

import (
	"github.com/rfashwall/kafka-deep-dive/pkg/utils"

	"github.com/rfashwall/kafka-deep-dive/cmd/inventory/consumer"

	log "github.com/sirupsen/logrus"

	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {

	//dbConn, err := db.NewDB()
	// if err != nil {
	// 	log.Panic("unable to connect to db.", err)
	// }
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

	c := consumer.Consumer{
		Broker: utils.BrokerAddress(),
		Group:  utils.ConsumerGroup(),
		Topic:  utils.OrderReceivedTopicName,
		//EventsDBManager: events.NewEventsDBManager(dbConn.PostgressDB),
		//DB:              dbConn,
	}

	log.Fatal(c.SubscribeAndListen())
}
