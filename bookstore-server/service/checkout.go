package service

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	dto "github.com/mtanzim/event-driven-bookstore/common-server/dto"
	"go.mongodb.org/mongo-driver/bson"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type CheckoutTopics struct {
	PaymentTopic  string
	ShipmentTopic string
}

type CheckoutService struct {
	producer   *kafka.Producer
	topics     *CheckoutTopics
	collection *mongo.Collection
}

func NewCheckoutService(p *kafka.Producer, topics *CheckoutTopics, collection *mongo.Collection) *CheckoutService {
	return &CheckoutService{p, topics, collection}
}

func (s CheckoutService) CheckoutCart(cart *dto.Cart) *dto.CartResponse {
	id := primitive.NewObjectID()
	go s.stageStock(cart.Items)
	go s.requestCartShipment(cart, id)
	go s.requestCartPayment(cart, id)
	cartResponse := dto.CartResponse{CartID: id, Status: "requested"}
	return &cartResponse

}

func (s CheckoutService) stageStock(items []dto.CartItem) {
	ctx, cancel := context.WithTimeout(context.Background(), CTXTimeout*time.Second)
	defer cancel()
	for _, item := range items {
		log.Println("updating stock for book", item.Book)
		log.Println("Following qty will be staged", item.Qty)
		filter := bson.M{"_id": item.Book.ID}
		update := bson.D{{"$inc", bson.M{"stagedQty": item.Qty}}}
		var updatedDocument bson.M
		err := s.collection.FindOneAndUpdate(ctx, filter, update, nil).Decode(&updatedDocument)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				log.Println("No documents found")
				return
			}
			log.Println(err)
			return
		}
		log.Printf("updated document %v", updatedDocument)
	}
}

// TODO: fire off shipment message
func (s CheckoutService) requestCartShipment(cart *dto.Cart, id primitive.ObjectID) {

	cartShipment := dto.CartShipment{
		CartID:  id,
		Address: cart.CartUserInformation.Address,
		Phone:   cart.CartUserInformation.Phone,
		Email:   cart.CartUserInformation.Email,
		Items:   cart.Items,
	}
	log.Println(cartShipment)
	topic := s.topics.ShipmentTopic
	msg, err := json.Marshal(cartShipment)
	if err != nil {
		log.Println(err)
	} else {
		s.sendMessageToProducer(&topic, msg)
	}

}

// TODO: add encryption to credit card data!
// https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes/
func (s CheckoutService) requestCartPayment(cart *dto.Cart, id primitive.ObjectID) {

	var totalPrice float64

	for _, item := range cart.Items {
		curPrice, err := strconv.ParseFloat(item.Book.Price, 64)
		if err != nil {
			curPrice = 0
		}
		totalPrice += curPrice * float64(item.Qty)
	}

	cartPayment := dto.CartPayment{
		CartID:     id,
		Address:    cart.CartUserInformation.Address,
		Phone:      cart.CartUserInformation.Phone,
		Email:      cart.CartUserInformation.Email,
		CardNum:    cart.CartUserInformation.CardNum,
		CardCode:   cart.CartUserInformation.CardCode,
		TotalPrice: totalPrice,
	}
	log.Println(cartPayment)

	topic := s.topics.PaymentTopic
	msg, err := json.Marshal(cartPayment)
	if err != nil {
		log.Println(err)
	} else {
		s.sendMessageToProducer(&topic, msg)
	}
}

func (s CheckoutService) sendMessageToProducer(topic *string, msg []byte) {

	err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
	log.Println(err)

}
