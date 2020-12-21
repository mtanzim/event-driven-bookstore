package main

import (
	"os"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	consumer "github.com/mtanzim/event-driven-bookstore/common-server/consumer"
	persister "github.com/mtanzim/event-driven-bookstore/common-server/persister"
	service "github.com/mtanzim/event-driven-bookstore/warehouse-server/service"
)

func main() {
	kafkaServer := os.Getenv("KAFKA_SERVER_ADDR")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	shipmentTopic := os.Getenv("SHIPMENT_TOPIC")
	paymentStatusTopic := os.Getenv("PROCESSED_PAYMENT_TOPIC")

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	db, disconnectDb := persister.NewMongo(uri, dbName)
	defer disconnectDb()
	collName := os.Getenv("MONGO_COLL")
	collection := db.Collection(collName)

	shipmentKafkaConsumer := consumer.NewKafkaConsumer(kafkaServer, groupID)
	shipmentService := service.NewShipmentService(shipmentKafkaConsumer, shipmentTopic, collection)

	paymentStatusConsumer := consumer.NewKafkaConsumer(kafkaServer, groupID)
	paymentStatusService := service.NewPaymentStatusService(paymentStatusConsumer, paymentStatusTopic, collection)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)
	go shipmentService.ConsumeMessages()
	go paymentStatusService.ConsumeMessages()
	wg.Wait()

}
