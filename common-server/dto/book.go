package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID     primitive.ObjectID ID      IDTag
	Title  string             `bson:"title" json:"title"`
	Author string             `bson:"author" json:"author"`
	Price  string             `bson:"price" json:"price"`
	Stock  int32              `bson:"stock" json:"stock"`
}
