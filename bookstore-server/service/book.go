package service

import (
	"context"
	"errors"
	"time"

	dto "github.com/mtanzim/event-driven-bookstore/common-server/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CTXTimeout = 5

type BookService struct {
	collection *mongo.Collection
}

func NewBookService(collection *mongo.Collection) *BookService {
	return &BookService{collection}
}

func (s BookService) GetBooks() ([]dto.Book, error) {
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
	bookResponse := generateBookResponse(dat)

	return bookResponse, nil

}

func generateBookResponse(dat []dto.Book) []dto.Book {
	var bookResponse []dto.Book
	for _, book := range dat {
		bookRes := book
		if newQty := book.Stock - book.StagedQty; newQty > 0 {
			bookRes.Stock = newQty
			bookResponse = append(bookResponse, bookRes)
		}
	}
	return bookResponse
}
