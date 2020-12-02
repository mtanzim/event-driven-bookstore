package service

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	collection *mongo.Collection
}

func NewService(db *mongo.Database, collName string) *Service {
	collection := db.Collection(collName)
	return &Service{collection}
}
