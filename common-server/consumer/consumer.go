package consumer

import (
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func NewKafkaConsumer(kafkaServerAddr string, groupId string) *kafka.Consumer {
	log.Println(kafkaServerAddr)
	log.Println(groupId)

	// connect to Kafka
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaServerAddr,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Printf("Consumer error: %v\n", err)
	}
	return c
}
