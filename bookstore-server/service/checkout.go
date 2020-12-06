package service

import (
	"encoding/json"
	"log"
	"strconv"

	dto "github.com/mtanzim/event-driven-bookstore/bookstore-server/dto"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func NewCheckoutService(p *kafka.Producer) *KafkaService {
	return &KafkaService{p}
}

func (s KafkaService) CheckoutCart(cart *dto.Cart) {
	id := primitive.NewObjectID()
	go s.requestCartShipment(cart, id)
	go s.requestCartPayment(cart, id)
}

// TODO: fire off shipment message
func (s KafkaService) requestCartShipment(cart *dto.Cart, id primitive.ObjectID) {

	cartShipment := dto.CartShipment{
		CartID:  id,
		Address: cart.CartUserInformation.Address,
		Phone:   cart.CartUserInformation.Phone,
		Items:   cart.Items,
	}
	log.Println(cartShipment)

}

// TODO: fire off a payment message
func (s KafkaService) requestCartPayment(cart *dto.Cart, id primitive.ObjectID) {

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
		TotalPrice: totalPrice,
	}
	log.Println(cartPayment)

	topic := "CartPayment"
	msg, err := json.Marshal(cartPayment)
	if err != nil {
		log.Println(err)
	} else {
		err = s.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          msg,
		}, nil)
		log.Println(err)

	}
}
