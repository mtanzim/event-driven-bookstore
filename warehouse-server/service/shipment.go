package service

import (
	"encoding/json"
	"log"

	"github.com/mtanzim/event-driven-bookstore/common-server/dto"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type ShipmentService struct {
	consumer   *kafka.Consumer
	topic      string
	collection *mongo.Collection
}

func NewShipmentService(c *kafka.Consumer, topic string, coll *mongo.Collection) *ShipmentService {
	return &ShipmentService{c, topic, coll}
}

func (s ShipmentService) ConsumeMessages() {
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
		var rawDat dto.CartShipment
		dat := dto.CartWarehouse{}
		shipment := json.Unmarshal()
		s.collection.InsertOne()
	}
}
