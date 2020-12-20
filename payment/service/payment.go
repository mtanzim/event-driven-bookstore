package service

import (
	"encoding/json"
	"log"

	"github.com/mtanzim/event-driven-bookstore/common-server/dto"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type PaymentService struct {
	consumer *kafka.Consumer
	topic    string
}

func NewPaymentService(c *kafka.Consumer, topic string) *PaymentService {
	return &PaymentService{c, topic}
}

func (s PaymentService) ConsumeMessages() {
	defer s.consumer.Close()
	log.Println(s.topic)
	s.consumer.Subscribe(s.topic, nil)
	for {
		msg, err := s.consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}
		log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

		var paymentRequest dto.CartPayment
		if err := json.Unmarshal(msg.Value, &paymentRequest); err != nil {
			log.Println(err)
		}
		go s.processPayment(&paymentRequest)

	}
}

func (s PaymentService) processPayment(paymentRequest *dto.CartPayment) {
	log.Println("Processing payment")
	log.Println(paymentRequest)
}
