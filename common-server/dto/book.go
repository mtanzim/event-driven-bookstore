package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Author    string             `bson:"author" json:"author"`
	Price     string             `bson:"price" json:"price"`
	Stock     int32              `bson:"stock" json:"stock"`
	StagedQty int32              `bson:"stagedQty" json:"stagedQty"`
}
