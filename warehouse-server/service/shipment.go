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
	warehouseService *WarehouseService
	collection       *mongo.Collection
}

func NewShipmentService(c *kafka.Consumer, topic string, coll *mongo.Collection) *ShipmentService {
	return &ShipmentService{&WarehouseService{c, topic}, coll}
}

func (s ShipmentService) GetPendingShipments() ([]dto.CartWarehouse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var data []dto.CartWarehouse
	cursor, err := s.collection.Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err

	}
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil

}

func (s ShipmentService) ConsumeMessages() {
	// TODO: is this idiomatic :/
	s.warehouseService.ConsumeMessages(s.processShipmentRequest)
}

func (s ShipmentService) processShipmentRequest(msg *kafka.Message) {

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
