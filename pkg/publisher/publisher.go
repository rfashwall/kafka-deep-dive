package publisher

import (
	"encoding/json"
	"pkg/pkg/events"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

func PublishEvent(event events.Event, topic string) error {
	log.WithField("event", event).Info("attempting to publish event")

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   "localhost",
		"socket.timeout.ms":   30000,
		"delivery.timeout.ms": 30000})

	if err != nil {
		return err
	}

	deliveryChan := make(chan kafka.Event)

	var value []byte

	if value, err = json.Marshal(event); err != nil {
		return err
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	log.WithField("Name", *m.TopicPartition.Topic).
		WithField("Partition", m.TopicPartition.Partition).
		WithField("PartitionOffset", m.TopicPartition.Offset).
		Infof("Delivered message to topic")

	close(deliveryChan)

	return nil
}