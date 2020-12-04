package service

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoService struct {
	collection *mongo.Collection
}

type KafkaService struct{}
