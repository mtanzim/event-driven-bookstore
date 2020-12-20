package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	consumer "github.com/mtanzim/event-driven-bookstore/common-server/consumer"
	service "github.com/mtanzim/event-driven-bookstore/payment/service"
)

func main() {
	kafkaServer := os.Getenv("KAFKA_SERVER_ADDR")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	paymentTopic := os.Getenv("PAYMENT_TOPIC")

	paymentKafkaConsumer := consumer.NewKafkaConsumer(kafkaServer, groupID)
	paymentService := service.NewPaymentService(paymentKafkaConsumer, paymentTopic)
	paymentService.ConsumeMessages()
}
