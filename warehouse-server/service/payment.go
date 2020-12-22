package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	dto "github.com/mtanzim/event-driven-bookstore/common-server/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type PaymentStatusService struct {
	warehouseService     *WarehouseService
	warehouseCollection  *mongo.Collection
	paymentDLQCollection *mongo.Collection
}

func NewPaymentStatusService(c *kafka.Consumer, topic string, warehouseColl, paymentDLQColl *mongo.Collection) *PaymentStatusService {
	return &PaymentStatusService{&WarehouseService{c, topic}, warehouseColl, paymentDLQColl}
}

func (s PaymentStatusService) ConsumeMessages() {
	s.warehouseService.ConsumeMessages(s.processPaymentStatus)
}

func (s PaymentStatusService) processPaymentStatus(msg *kafka.Message) {

	var paymentStatus dto.CartPaymentResponse
	if err := json.Unmarshal(msg.Value, &paymentStatus); err != nil {
		log.Println(err)
	}
	go s.persist(paymentStatus)
}

type CartPaymentDLQItem struct {
	CartID primitive.ObjectID `bson:"cartId,omitempty" json:"cartId"`
}

// TODO: synchronize? What if payment notification comes before the shipment request is registered
func (s PaymentStatusService) persist(paymentStatus dto.CartPaymentResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("Updating payment status")
	log.Println(paymentStatus)

	if paymentStatus.Approved {
		update := bson.D{{"$set", bson.D{{"paid", true}}}}
		filter := bson.M{"_id": paymentStatus.CartID}
		updateRes, err := s.warehouseCollection.UpdateOne(ctx, filter, update)
		if err != nil || updateRes.ModifiedCount != 1 {
			log.Println("Failed to update shipment payment status for cart id:", paymentStatus.CartID)
			log.Println("Inserting approved payment into DLQ")
			// TODO: can this be done with Kafka? Come back to this later
			failedCartItem := CartPaymentDLQItem{CartID: paymentStatus.CartID}
			s.paymentDLQCollection.InsertOne(ctx, failedCartItem)
		} else {
			log.Println("Shipment", paymentStatus.CartID, "was paid for!")
		}

	}

}
