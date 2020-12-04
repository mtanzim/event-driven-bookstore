package service

import (
	"log"
	"strconv"

	dto "github.com/mtanzim/event-driven-bookstore/bookstore-server/dto"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

func NewCheckoutService() *KafkaService {
	return &KafkaService{}
}

func (s KafkaService) CheckoutCart(cart *dto.Cart) error {
	id := primitive.NewObjectID()
	go s.requestCartShipment(cart, id)
	go s.requestCartPayment(cart, id)
	return nil
}

// TODO: fire off shipment message
func (s KafkaService) requestCartShipment(cart *dto.Cart, id primitive.ObjectID) error {

	cartShipment := dto.CartShipment{
		CartID:  id,
		Address: cart.CartUserInformation.Address,
		Phone:   cart.CartUserInformation.Phone,
		Items:   cart.Items,
	}
	log.Println(cartShipment)

	return nil

}

// TODO: fire off a payment message
func (s KafkaService) requestCartPayment(cart *dto.Cart, id primitive.ObjectID) error {

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

	return nil

}
