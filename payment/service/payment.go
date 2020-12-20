package service

import (
	"encoding/json"
	"log"

	"github.com/mtanzim/event-driven-bookstore/common-server/dto"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type PaymentService struct {
	consumer              *kafka.Consumer
	producer              *kafka.Producer
	topicPaymentRequested string
	topicPaymentProcessed string
}

func NewPaymentService(consumer *kafka.Consumer, topicPaymentRequested string, producer *kafka.Producer, topicPaymentProcessed string) *PaymentService {
	return &PaymentService{consumer, producer, topicPaymentRequested, topicPaymentProcessed}
}

func (s PaymentService) ConsumeMessages() {
	defer s.consumer.Close()
	log.Println(s.topicPaymentRequested)
	s.consumer.Subscribe(s.topicPaymentRequested, nil)
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

// TODO: 3rd party payment processing?
func (s PaymentService) processPayment(paymentRequest *dto.CartPayment) {
	log.Println("Processing payment")
	log.Println(paymentRequest)

	paymentStatus := dto.CartPaymentResponse{CartID: paymentRequest.CartID,
		Approved: true,
	}

	topic := s.topicPaymentProcessed
	msg, err := json.Marshal(paymentStatus)
	if err != nil {
		log.Println(err)
	} else {
		s.sendMessageToProducer(&topic, msg)
	}
}

// TODO: move this to common?
func (s PaymentService) sendMessageToProducer(topic *string, msg []byte) {

	err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
	log.Println(err)

}
