package service

import (
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type WarehouseService struct {
	consumer *kafka.Consumer
	topic    string
}

func (s WarehouseService) ConsumeMessages(processingFn func(msg *kafka.Message)) {
	defer s.consumer.Close()
	log.Println(s.topic)
	s.consumer.Subscribe(s.topic, nil)
	for {
		msg, err := s.consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}
		log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		go processingFn(msg)
	}
}
