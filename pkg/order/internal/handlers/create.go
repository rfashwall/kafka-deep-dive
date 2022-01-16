package handlers

import (
	"encoding/json"
	"net/http"
	events "pkg/pkg/models"
	ordermodel "pkg/pkg/models/order"
	"pkg/pkg/publisher"
	"pkg/pkg/utils"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func Create(w http.ResponseWriter, r *http.Request) {
	var order ordermodel.Order

	order.ID = uuid.New()
	var err error

	if err = json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.WithField("order", order).Info("received new order")
	orderEvent := orderToEvent(order)
	log.WithField("event", orderEvent).Info("transformed order to event")

	if err = publisher.PublishEvent(orderEvent, utils.OrderReceivedTopicName); err != nil {

		log.WithField("orderID", order.ID).Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithField("event", orderEvent).Info("published order event")
	w.WriteHeader(http.StatusCreated)
}

func orderToEvent(order ordermodel.Order) events.Event {
	return ordermodel.OrderReceived{
		EventBase: events.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
		},
		EventBody: order,
	}
}
