package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostShipment struct {
	CartID primitive.ObjectID `bson:"_id,omitempty" json:"cartId"`
}
