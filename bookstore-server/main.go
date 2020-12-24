package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mtanzim/event-driven-bookstore/bookstore-server/handler"
	"github.com/mtanzim/event-driven-bookstore/bookstore-server/service"
	"github.com/mtanzim/event-driven-bookstore/common-server/persister"
	kafkaProducer "github.com/mtanzim/event-driven-bookstore/common-server/producer"
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

	var checkoutTopics service.CheckoutTopics
	checkoutTopics.PaymentTopic = os.Getenv("PAYMENT_TOPIC")
	checkoutTopics.ShipmentTopic = os.Getenv("SHIPMENT_TOPIC")

	db, disconnectDb := persister.NewMongo(uri, dbName)
	defer disconnectDb()

	kafkaServerAddr := os.Getenv("KAFKA_SERVER_ADDR")
	producer := kafkaProducer.NewKafkaProducer(kafkaServerAddr)
	defer producer.Close()

	bookCollection := db.Collection("books")

	bookHandler := handler.NewBookHandler(bookCollection)
	checkoutHandler := handler.NewCheckoutHandler(producer, &checkoutTopics, bookCollection)
	r := mux.NewRouter()
	port := os.Getenv("REST_PORT")
	r.HandleFunc("/api/books", bookHandler.GetBooks).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/checkout", checkoutHandler.Checkout).Methods(http.MethodPost, http.MethodOptions)
	r.Use(loggingMiddleware)
	// TODO: fix this
	rWithCORS := cors.Default().Handler(r)
	log.Println("Starting HTTP Server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, rWithCORS))
}
