package consumer

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func NewKafkaConsumer(kafkaServerAddr string, groupId string) *kafka.Consumer {
	// connect to Kafka
	log.Println(kafkaServerAddr)
	log.Println(groupId)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaServerAddr,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Printf("Consumer error: %v\n", err)
	}
	return c
}

// func ConsumeShipmentMessages(c *kafka.Consumer, topic string) {
// 	defer c.Close()
// 	log.Println(topic)
// 	c.Subscribe(topic, nil)
// 	consumeMessages(c)
// }

// func ConsumePaymentMessages() {
// 	defer c.Close()
// 	log.Println(topic)
// 	c.Subscribe(topic, nil)
// 	consumeMessages(c)

// }

func ConsumeMessages(c *kafka.Consumer, topic string) {
	defer c.Close()
	log.Println(topic)
	c.Subscribe(topic, nil)
	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}
		log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
	}
}
