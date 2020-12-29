package handler

import (
	"encoding/json"
	"log"
	"net/http"

	localDTO "github.com/mtanzim/event-driven-bookstore/warehouse-server/dto"
	service "github.com/mtanzim/event-driven-bookstore/warehouse-server/service"
)

type ShipmentHandler struct {
	service *service.ShipmentService
}

func NewShipmentHandler(s *service.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{s}
}

func (h ShipmentHandler) GetPendingShipments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pendingShipments, err := h.service.GetPendingShipments()
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot get data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pendingShipments)
}

func (h ShipmentHandler) PostPendingShipemt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cart localDTO.PostShipment
	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		log.Println(err)
		http.Error(w, "Cannot process shipment", http.StatusBadRequest)
		return
	}
	postedShipment, err := h.service.PostPendingShipemt(&cart)
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot update shipment", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(postedShipment)
}
