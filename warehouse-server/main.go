package main

import (
	"os"
	"sync"

	consumer "github.com/mtanzim/event-driven-bookstore/warehouse-server/consumer"
)

func main() {
	kafkaServer := os.Getenv("KAFKA_SERVER_ADDR")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	shipmentTopic := os.Getenv("SHIPMENT_TOPIC")
	paymentTopic := os.Getenv("PAYMENT_TOPIC")

	shipmentConsumer := consumer.NewKafkaConsumer(kafkaServer, groupID)
	paymentConsumer := consumer.NewKafkaConsumer(kafkaServer, groupID)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)
	go consumer.ConsumeMessages(paymentConsumer, paymentTopic, &wg)
	go consumer.ConsumeMessages(shipmentConsumer, shipmentTopic, &wg)
	wg.Wait()

}
