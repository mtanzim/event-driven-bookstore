package service

import "go.mongodb.org/mongo-driver/bson/primitive"

type CartPaymentDLQItem struct {
	CartID primitive.ObjectID `bson:"cartId,omitempty" json:"cartId"`
}
