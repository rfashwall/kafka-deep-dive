package main

import (
	"os"
	events "pkg/pkg/models"
	order "pkg/pkg/models/order"
	"pkg/pkg/publisher"
	"pkg/pkg/utils"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	var err error
	var o = order.Order{
		ID: uuid.New(),
	}

	var event = order.OrderReceived{
		EventBase: events.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
		},
		EventBody: o,
	}

	if err = publisher.PublishEvent(event, utils.OrderReceivedTopicName); err != nil {
		log.WithField("orderID", o.ID).Error(err.Error())
	} else {
		log.WithField("event", event).Info("published event")
	}
}
