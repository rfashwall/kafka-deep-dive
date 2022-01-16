package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/rfashwall/kafka-deep-dive/pkg/db"
	eventsdb "github.com/rfashwall/kafka-deep-dive/pkg/events"
	events "github.com/rfashwall/kafka-deep-dive/pkg/models"
	order "github.com/rfashwall/kafka-deep-dive/pkg/models/order"
	"github.com/rfashwall/kafka-deep-dive/pkg/publisher"
	"github.com/rfashwall/kafka-deep-dive/pkg/utils"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	Broker          string
	Group           string
	Topic           string
	EventsDBManager eventsdb.DBManager
	DB              db.DB
}

func (c *Consumer) SubscribeAndListen() error {

	kc, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     c.Broker,
		"broker.address.family": "v4",
		"group.id":              c.Group + "-inventory",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	})

	if err != nil {
		log.WithField("error", err).Error("Failed to create consumer")

		return err
	}
	log.WithField("consumer", kc).Info("Created Consumer")

	if err = kc.SubscribeTopics([]string{c.Topic}, nil); err != nil {
		log.WithField("error", err).
			WithField("topic", c.Topic).
			Error("Failed to subscribe to topic")

		return err
	}

	for {
		msg, err := kc.ReadMessage(-1)
		if err != nil {
			// The client will automatically try to recover from all errors.
			log.WithField("error", err).Error(msg)

			log.Warn("Closing consumer...")
			kc.Close()

			return err
		}

		log.WithField("topic", msg.TopicPartition).Info(string(msg.Value))

		var event order.OrderReceived
		if err = json.Unmarshal([]byte(string(msg.Value)), &event); err != nil {
			log.WithField("error", err).Error("an issue occurred unmarshalling event from message received")

			continue
		}

		var order order.Order
		if order, err = extractOrder(event); err != nil {
			log.WithField("error", err).Error("an issue occurred trying to extract order information from the order recieved event")

			//hdlr.HandleError(event)
			continue
		}

		if err = c.processEvent(event, order); err != nil {
			log.WithField("error", err).Error("an issue occurred trying to process the event")

			//hdlr.HandleError(event)
			continue
		}
	}
}

func extractOrder(event order.OrderReceived) (order.Order, error) {
	log.Info("attempting to extract order from event")

	body := event.Body()
	o, ok := body.(order.Order)
	if !ok {
		return order.Order{}, errors.New("event body can't be cast as an order")
	}

	return o, nil
}

func (c *Consumer) processEvent(event events.Event, order order.Order) error {
	publishOrderConfirmedEvent(order)
	tx, err := c.DB.PostgressDB.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	c.EventsDBManager.EventExists(event, tx)
	// check to see if event has already been processed
	eventAlreadyProcessed, err := c.EventsDBManager.EventExists(event, tx)
	if err != nil {
		log.WithField("error", err).Error("an issue occurred trying to check if an event was already processed")
		return err
	}

	// if event has already been processed, nothing more to do
	if eventAlreadyProcessed {
		log.WithField("event.id", event.ID()).
			WithField("event.name", event.Name()).
			Info("event was processed previously")

		return nil
	}

	// mark the event as processed
	if err = c.EventsDBManager.InsertEvent(event, tx); err != nil {
		log.WithField("error", err).Error("an issue occurred trying to insert the event")
		return err
	}

	return nil
}

func publishOrderConfirmedEvent(o order.Order) error {
	// publish an order confirmed event
	e := translateOrderToEvent(o)

	log.WithField("event", e).Info("transformed order to event")

	var err error
	if err = publisher.PublishEvent(e, utils.OrderConfirmedTopicName); err != nil {
		return err
	}

	log.WithField("event", e).Info("published event")

	return nil
}

func translateOrderToEvent(o order.Order) events.Event {
	var event = order.OrderConfirmed{
		EventBase: events.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
		},
		EventBody: o,
	}

	return event
}
