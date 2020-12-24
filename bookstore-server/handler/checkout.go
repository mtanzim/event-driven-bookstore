package handler

import (
	"encoding/json"
	"log"
	"net/http"

	service "github.com/mtanzim/event-driven-bookstore/bookstore-server/service"
	dto "github.com/mtanzim/event-driven-bookstore/common-server/dto"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type CheckoutHandler struct {
	service *service.CheckoutService
}

func NewCheckoutHandler(p *kafka.Producer, topics *service.CheckoutTopics, collection *mongo.Collection) *CheckoutHandler {
	s := service.NewCheckoutService(p, topics, collection)
	return &CheckoutHandler{s}
}

// TODO: form validations
func (h CheckoutHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text")
	var cart dto.Cart
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot post cart", http.StatusBadRequest)
		return
	}
	cartRes := h.service.CheckoutCart(&cart)
	json.NewEncoder(w).Encode(cartRes)
}
