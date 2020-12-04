package handler

import (
	"encoding/json"
	"log"
	"net/http"

	dto "github.com/mtanzim/event-driven-bookstore/bookstore-server/dto"
	service "github.com/mtanzim/event-driven-bookstore/bookstore-server/service"
)

type CheckoutHandler struct {
	service *service.KafkaService
}

func NewCheckoutHandler() *CheckoutHandler {
	s := service.NewCheckoutService()
	return &CheckoutHandler{s}
}

func (h CheckoutHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text")
	var cart dto.Cart
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot post cart", http.StatusBadRequest)
		return
	}
	go h.service.CheckoutCart(&cart)
	w.Write([]byte("OK"))
}
