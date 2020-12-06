package service

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type MongoService struct {
	collection *mongo.Collection
}

type KafkaService struct {
	producer *kafka.Producer
}
