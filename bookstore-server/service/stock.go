package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/mtanzim/event-driven-bookstore/common-server/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// TODO: does it make sense to have stock in bookstore?
type StockUpdateService struct {
	consumer   *kafka.Consumer
	topic      string
	collection *mongo.Collection
}

func NewStockUpdateService(consumer *kafka.Consumer, topic string, coll *mongo.Collection) *StockUpdateService {
	return &StockUpdateService{consumer, topic, coll}
}

func (s StockUpdateService) updateStock(msg *kafka.Message) {
	var shippedCart dto.CartWarehouseShipped
	if err := json.Unmarshal(msg.Value, &shippedCart); err != nil {
		log.Println(err)
	}
	go s.persist(shippedCart)

}

func (s StockUpdateService) persist(paymentStatus dto.CartWarehouseShipped) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("Updating stock")
	log.Println(paymentStatus)
	items := paymentStatus.Cart.Items
	for _, item := range items {
		log.Println("Updating book stock")
		log.Println(item.Book)
		log.Println(item.Qty)
		filter := bson.M{"_id": item.Book.ID}
		update := bson.D{{"$inc", bson.D{{"stagedQty", -item.Qty}}}, {"$inc", bson.D{{"stock", -item.Qty}}}}
		updateRes, err := s.collection.UpdateOne(ctx, filter, update)
		if err != nil || updateRes.ModifiedCount != 1 {
			log.Println("Failed to update stock for book id:", item.Book.ID)
			log.Println(err)
		} else {
			log.Println("Book", item.Book.ID, "stock was updated")
		}

	}
}

func (s StockUpdateService) MonitorAndUpdate() {
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
		s.updateStock(msg)
	}

}
