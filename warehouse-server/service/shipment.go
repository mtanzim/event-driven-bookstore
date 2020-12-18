package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	dto "github.com/mtanzim/event-driven-bookstore/common-server/dto"
	"go.mongodb.org/mongo-driver/bson"
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
		var cart dto.CartShipment
		if err := json.Unmarshal(msg.Value, &cart); err != nil {
			log.Println(err)
		}

		shipmentRequest := dto.CartWarehouse{
			ID:      cart.CartID,
			Cart:    cart,
			Shipped: false,
			Paid:    false,
		}
		log.Println(shipmentRequest)
		go s.persist(shipmentRequest)
	}
}

func (s ShipmentService) persist(shipmentRequest dto.CartWarehouse) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing bson.M
	if err := s.collection.FindOne(ctx, bson.M{"_id": shipmentRequest.ID}).Decode(&existing); err != nil {
		if err == mongo.ErrNoDocuments {
			insertRes, insertErr := s.collection.InsertOne(ctx, shipmentRequest)
			if insertErr != nil {
				// TODO: ignore failed insert on duplicate?
				log.Println(err)
			}
			log.Println("Inserted shipment ", insertRes.InsertedID)
		}
	}

}
