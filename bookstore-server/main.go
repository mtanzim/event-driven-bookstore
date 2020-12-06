package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	broker "github.com/mtanzim/event-driven-bookstore/bookstore-server/broker"
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

func main() {

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	db, disconnectDb := persister.NewMongo(uri, dbName)
	defer disconnectDb()

	kafkaServerAddr := os.Getenv("KAFKA_SERVER_ADDR")
	producer := broker.NewKafkaProducer(kafkaServerAddr)
	defer producer.Close()

	bookHandler := handler.NewBookHandler(db)
	checkoutHandler := handler.NewCheckoutHandler(producer)
	r := mux.NewRouter()
	port := os.Getenv("REST_PORT")
	r.HandleFunc("/api/books", bookHandler.GetBooks).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/checkout", checkoutHandler.Checkout).Methods(http.MethodPost, http.MethodOptions)
	r.Use(loggingMiddleware)
	// TODO: fix this
	rWithCORS := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":"+port, rWithCORS))
}
