package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	dto "github.com/mtanzim/event-driven-bookstore/common-server/dto"
	localDTO "github.com/mtanzim/event-driven-bookstore/warehouse-server/dto"
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
	filter := bson.M{"shipped": false}
	cursor, err := s.collection.Find(ctx, filter, nil)
	if err != nil {
		return nil, err

	}
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil

}

func (s ShipmentService) PostPendingShipemt(cart *localDTO.PostShipment) (*localDTO.PostShipment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": cart.CartID}
	update := bson.D{{"$set", bson.D{{"shipped", true}}}}

	updateRes, updateErr := s.collection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		return nil, updateErr
	}
	if updateRes.ModifiedCount == 1 || updateRes.UpsertedCount == 1 {
		log.Println("Successfully shipped cart:", cart.CartID)
		return cart, nil
	}
	return nil, errors.New("Something went wrong.")

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
