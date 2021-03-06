package handler

import (
	"encoding/json"
	"log"
	"net/http"

	service "github.com/mtanzim/event-driven-bookstore/bookstore-server/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookHandler struct {
	service *service.BookService
}

func NewBookHandler(collection *mongo.Collection) *BookHandler {
	s := service.NewBookService(collection)
	return &BookHandler{s}
}

func (h BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books, err := h.service.GetBooks()
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot get data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}
