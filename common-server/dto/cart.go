package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// TODO: how to implement DRY here?

type CartItem struct {
	Book BookInCart `json:"book"`
	Qty  int32      `json:"qty"`
}

type CartUserInformation struct {
	Address  string `json:"address"`
	CardNum  string `json:"cardNum"`
	CardCode string `json:"code"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type Cart struct {
	Items               []CartItem          `json:"items"`
	CartUserInformation CartUserInformation `json:"cartUserInformation`
}

type CartPayment struct {
	CartID     primitive.ObjectID `bson:"_id,omitempty" json:"cartId"`
	Address    string             `json:"address"`
	CardNum    string             `json:"cardNum"`
	CardCode   string             `json:"code"`
	Email      string             `json:"email"`
	Phone      string             `json:"phone"`
	TotalPrice float64            `json:"totalPrice"`
}

type CartShipment struct {
	CartID  primitive.ObjectID `bson:"_id,omitempty" json:"cartId"`
	Address string             `json:"address"`
	Phone   string             `json:"phone"`
	Email   string             `json:"email"`
	Items   []CartItem         `json:"items"`
}

type CartResponse struct {
	CartID primitive.ObjectID `bson:"_id,omitempty" json:"cartId"`
	Status string             `json:"status"`
}

type CartPaymentResponse struct {
	CartID   primitive.ObjectID `bson:"_id,omitempty" json:"cartId"`
	Approved bool               `json:"approved"`
	Message  string             `json:"message,omitempty"`
}

type CartWarehouse struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"cartId"`
	Cart    CartShipment       `json:"cart"`
	Shipped bool               `json:"shipped" bson:"shipped"`
	Paid    bool               `json:"paid" bson:"paid"`
}

type CartWarehouseShipped struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"cartId"`
	Cart CartShipment       `json:"cart"`
}
