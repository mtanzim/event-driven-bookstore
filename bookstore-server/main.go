package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	handler "github.com/mtanzim/event-driven-bookstore/bookstore-server/handler"
	persister "github.com/mtanzim/event-driven-bookstore/bookstore-server/persister"
	kafkaProducer "github.com/mtanzim/event-driven-bookstore/bookstore-server/producer"
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

	topics := make(map[string]string)
	topics["CART_TOPIC"] = os.Getenv("CART_TOPIC")
	topics["SHIPMENT_TOPIC"] = os.Getenv("SHIPMENT_TOPIC")

	db, disconnectDb := persister.NewMongo(uri, dbName)
	defer disconnectDb()

	kafkaServerAddr := os.Getenv("KAFKA_SERVER_ADDR")
	producer := kafkaProducer.NewKafkaProducer(kafkaServerAddr)
	defer producer.Close()

	bookHandler := handler.NewBookHandler(db)
	checkoutHandler := handler.NewCheckoutHandler(producer, topics)
	r := mux.NewRouter()
	port := os.Getenv("REST_PORT")
	r.HandleFunc("/api/books", bookHandler.GetBooks).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/checkout", checkoutHandler.Checkout).Methods(http.MethodPost, http.MethodOptions)
	r.Use(loggingMiddleware)
	// TODO: fix this
	rWithCORS := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":"+port, rWithCORS))
}
