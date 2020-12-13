package main

import (
	"log"
	"os"

	consumer "github.com/mtanzim/event-driven-bookstore/warehouse-server/consumer"
)

func main() {
	log.Println("Hello World")
	kafkaServer := os.Getenv("KAFKA_SERVER_ADDR")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	// shipmentTopic := os.Getenv("SHIPMENT_TOPIC")
	paymentTopic := os.Getenv("PAYMENT_TOPIC")
	// shipmentConsumer := consumer.NewKafkaConsumer(kafkaServer, groupID)
	paymentConsumer := consumer.NewKafkaConsumer(kafkaServer, groupID)
	// go consumer.ConsumeMessages(shipmentConsumer, shipmentTopic)
	consumer.ConsumeMessages(paymentConsumer, paymentTopic)

}
