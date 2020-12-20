package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	consumer "github.com/mtanzim/event-driven-bookstore/common-server/consumer"
	producer "github.com/mtanzim/event-driven-bookstore/common-server/producer"
	service "github.com/mtanzim/event-driven-bookstore/payment/service"
)

func main() {
	kafkaServer := os.Getenv("KAFKA_SERVER_ADDR")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	paymentTopic := os.Getenv("PAYMENT_TOPIC")
	paymentProcessedTopic := os.Getenv("PROCESSED_PAYMENT_TOPIC")

	paymentKafkaConsumer := consumer.NewKafkaConsumer(kafkaServer, groupID)
	paymentKafkaProducer := producer.NewKafkaProducer(kafkaServer)
	defer paymentKafkaProducer.Close()

	paymentService := service.NewPaymentService(paymentKafkaConsumer, paymentTopic, paymentKafkaProducer, paymentProcessedTopic)
	paymentService.ConsumeMessages()
}
