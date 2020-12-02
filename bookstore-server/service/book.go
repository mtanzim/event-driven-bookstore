package service

import (
	"context"
	"errors"
	"time"

	dto "github.com/mtanzim/event-driven-bookstore/bookstore-server/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CTXTimeout = 5
const CollectionName = "books"

func NewBookService(db *mongo.Database) *Service {
	return NewService(db, CollectionName)
}

func (s Service) GetBooks() ([]dto.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CTXTimeout*time.Second)
	defer cancel()

	opts := options.Find()
	opts.SetSort(bson.D{{"author", 1}})
	var dat []dto.Book
	cursor, err := s.collection.Find(ctx, bson.M{}, opts)
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &dat); err != nil {
		return nil, errors.New("Cannot get data")
	}

	return dat, nil

}
