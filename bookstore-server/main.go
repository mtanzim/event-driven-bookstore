package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	dto "github.com/mtanzim/event-driven-bookstore/bookstore-server/dto"
	handler "github.com/mtanzim/event-driven-bookstore/bookstore-server/handler"
	persister "github.com/mtanzim/event-driven-bookstore/bookstore-server/persister"
	"github.com/rs/cors"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func checkout(w http.ResponseWriter, r *http.Request) {

	var cart []dto.CartItem
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		http.Error(w, "Cannot post cart", http.StatusInternalServerError)

	}
	log.Println(cart)

}

func main() {
	db, disconnectDb := persister.NewMongo()
	defer disconnectDb()
	bookHandler := handler.NewBookHandler(db)
	r := mux.NewRouter()
	port := os.Getenv("REST_PORT")
	r.HandleFunc("/api/books", bookHandler.GetBooks).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/checkout", checkout).Methods(http.MethodPost, http.MethodOptions)
	r.Use(loggingMiddleware)
	// TODO: fix this
	rWithCORS := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":"+port, rWithCORS))
}
