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

// TODO: synchronize?
// Is this a proper way to build a DLQ?
// Can this be done with Kafka? Come back to this later.
func (s PaymentStatusService) persist(paymentStatus dto.CartPaymentResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("Updating payment status")
	log.Println(paymentStatus)

	if paymentStatus.Approved {
		filter := bson.M{"_id": paymentStatus.CartID}
		update := bson.D{{"$set", bson.D{{"paid", true}}}}
		updateRes, err := s.warehouseCollection.UpdateOne(ctx, filter, update)
		if err != nil || updateRes.ModifiedCount != 1 {
			log.Println("Failed to update shipment payment status for cart id:", paymentStatus.CartID)
			log.Println("Inserting approved payment into DLQ")
			failedCartItem := CartPaymentDLQItem{CartID: paymentStatus.CartID}
			s.paymentDLQCollection.InsertOne(ctx, failedCartItem)
		} else {
			log.Println("Shipment", paymentStatus.CartID, "was paid for!")
		}

	}
	// TODO: Payment rejection logic
	// if rejected, fire of 'unstage' message on broker with []CartItem
	// store microservice can update staged and stock qty
	// Alternatively store can directly listen for payment rejections?
	// Then, store will need cart persistence as well

}
