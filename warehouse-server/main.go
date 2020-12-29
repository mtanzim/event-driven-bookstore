package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	consumer "github.com/mtanzim/event-driven-bookstore/common-server/consumer"
	persister "github.com/mtanzim/event-driven-bookstore/common-server/persister"
	"github.com/mtanzim/event-driven-bookstore/common-server/producer"
	handler "github.com/mtanzim/event-driven-bookstore/warehouse-server/handler"
	service "github.com/mtanzim/event-driven-bookstore/warehouse-server/service"
	"github.com/rs/cors"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	kafkaServer := os.Getenv("KAFKA_SERVER_ADDR")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	shipmentTopic := os.Getenv("SHIPMENT_TOPIC")
	paymentStatusTopic := os.Getenv("PROCESSED_PAYMENT_TOPIC")
	shipmentCompletedTopic := os.Getenv("CART_SHIPMENT_COMPLETE_TOPIC")

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	db, disconnectDb := persister.NewMongo(uri, dbName)
	defer disconnectDb()
	warehouseCollName := os.Getenv("MONGO_COLL_WAREHOUSE")
	warehouseCollection := db.Collection(warehouseCollName)

	warehousePaymentDLQCollName := os.Getenv("MONGO_COLL_PAYMENT_PENDING_ACK")
	warehousePaymentDLQColl := db.Collection(warehousePaymentDLQCollName)

	shipmentKafkaConsumer := consumer.NewKafkaConsumer(kafkaServer, groupID)
	shipmentKafkaProducer := producer.NewKafkaProducer(kafkaServer)
	// TODO: arg list too long! Combine into struct?
	shipmentService := service.NewShipmentService(shipmentKafkaConsumer, shipmentKafkaProducer, shipmentTopic, shipmentCompletedTopic, warehouseCollection)

	paymentStatusConsumer := consumer.NewKafkaConsumer(kafkaServer, groupID)
	paymentStatusService := service.NewPaymentStatusService(paymentStatusConsumer, paymentStatusTopic, warehouseCollection, warehousePaymentDLQColl)

	dqlMonitorService := service.NewDLQMonitorService(warehousePaymentDLQColl, warehouseCollection)

	go dqlMonitorService.Monitor()
	go shipmentService.ConsumeMessages()
	go paymentStatusService.ConsumeMessages()

	shipmentHandler := handler.NewShipmentHandler(shipmentService)
	r := mux.NewRouter()
	port := os.Getenv("REST_PORT")
	r.HandleFunc("/api/shipment", shipmentHandler.GetPendingShipments).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/shipment", shipmentHandler.PostPendingShipemt).Methods(http.MethodPost, http.MethodOptions)
	r.Use(loggingMiddleware)
	// TODO: fix this
	rWithCORS := cors.Default().Handler(r)
	log.Println("Starting HTTP Server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, rWithCORS))

}
