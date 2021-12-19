package main

import (
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type orderHandlers struct {
}

func (h *orderHandlers) order(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	bt := []byte("ok")
	w.Write(bt)
}

func main() {
	go initPublisher()
	orderHandlers := neworderHandlers()
	http.HandleFunc("/order", orderHandlers.order)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}

func neworderHandlers() *orderHandlers {
	return &orderHandlers{}
}

func initPublisher() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "OrderReceived"
	jsonObj := `
	{
		"EventBase": {
		  "EventID": "1774f854-07a2-4b85-be12-db31d1684bba",
		  "EventName": "OrderReceived",
		  "EventTimestamp": "2020-08-16T17:47:39.087555-04:00"
		},
		"EventBody": {
		  "id": "6e042f29-350b-4d51-8849-5e36456dfa48",
		  "products": [
			{
			  "productCode": "12345",
			  "quantity": 2
			}
		  ],
		  "customer": {
			"firstName": "Tom",
			"lastName": "Hardy",
			"emailAddress": "tom.hardy@email.com",
			"shippingAddress": {
			  "line1": "123 Anywhere St",
			  "city": "Anytown",
			  "state": "AL",
			  "postalCode": "12345"
			}
		  }
		}
	  }`
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(jsonObj),
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
