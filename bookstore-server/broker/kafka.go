package broker

import (
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func kafkaDeliveryReports(p *kafka.Producer) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				log.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		default:
			log.Println(ev)
		}
	}
}

func NewKafkaProducer(kafkaServerAddr string) *kafka.Producer {
	log.Println(kafkaServerAddr)
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaServerAddr})
	if err != nil {
		log.Fatalln(err)
	}
	go kafkaDeliveryReports(p)
	return p
}
