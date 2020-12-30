package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	dto "github.com/mtanzim/event-driven-bookstore/common-server/dto"
	localDTO "github.com/mtanzim/event-driven-bookstore/warehouse-server/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type ShipmentService struct {
	warehouseService       *WarehouseService
	collection             *mongo.Collection
	producer               *kafka.Producer
	shipmentCompletedTopic string
}

func NewShipmentService(c *kafka.Consumer, p *kafka.Producer, shipmentTopic, shipmentCompletedTopic string, coll *mongo.Collection) *ShipmentService {
	return &ShipmentService{&WarehouseService{c, shipmentTopic}, coll, p, shipmentCompletedTopic}
}

func (s ShipmentService) GetPendingShipments() ([]dto.CartWarehouse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var data []dto.CartWarehouse
	filter := bson.M{"shipped": false}
	cursor, err := s.collection.Find(ctx, filter, nil)
	if err != nil {
		return nil, err

	}
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}

	if data == nil {
		return []dto.CartWarehouse{}, nil
	}
	return data, nil

}

func (s ShipmentService) PostPendingShipemt(cart *localDTO.PostShipment) (*dto.CartWarehouseShipped, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": cart.CartID}
	update := bson.D{{"$set", bson.D{{"shipped", true}}}}

	var cartToUpdate dto.CartWarehouseShipped
	updateErr := s.collection.FindOneAndUpdate(ctx, filter, update).Decode(&cartToUpdate)
	if updateErr != nil {
		return nil, updateErr
	}
	go s.sendMessageToProducer(&cartToUpdate)
	return &cartToUpdate, nil

}

func (s ShipmentService) sendMessageToProducer(cart *dto.CartWarehouseShipped) {
	msg, err := json.Marshal(cart)
	if err != nil {
		log.Println(err)
		return
	}
	err = s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.shipmentCompletedTopic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
	log.Println(err)

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
