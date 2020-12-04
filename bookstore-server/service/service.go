package service

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoService struct {
	collection *mongo.Collection
}

type SimpleService struct{}

func NewMongoService(db *mongo.Database, collName string) *MongoService {
	collection := db.Collection(collName)
	return &MongoService{collection}
}

func NewSimpleService() *SimpleService {
	return &SimpleService{}
}
