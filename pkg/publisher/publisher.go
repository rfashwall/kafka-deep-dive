package publisher

import (
	"encoding/json"

	"github.com/rfashwall/kafka-deep-dive/pkg/models"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

func PublishEvent(event models.Event, topic string) error {
	log.WithField("event", event).Info("attempting to publish event")

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   "localhost",
		"socket.timeout.ms":   30000,
		"delivery.timeout.ms": 30000})

	if err != nil {
		return err
	}

	deliveryChan := make(chan kafka.Event)

	value, err := json.Marshal(event)

	if err != nil {
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
