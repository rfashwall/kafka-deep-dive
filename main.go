package main

import (
	"os"
	"pkg/pkg/events"
	orderevents "pkg/pkg/events/order"
	ordermodel "pkg/pkg/order/pkg/model"
	"pkg/pkg/publisher"
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
	log.SetLevel(0)
}

func main() {
	var err error
	var order = ordermodel.Order{
		ID: uuid.New(),
	}

	var event = orderevents.OrderReceived{
		EventBase: events.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
		},
		EventBody: order,
	}

	if err = publisher.PublishEvent(event, "OrderReceived"); err != nil {
		log.WithField("orderID", order.ID).Error(err.Error())
	} else {
		log.WithField("event", event).Info("published event")
	}
}
