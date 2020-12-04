package dto

type CartItem struct {
	Book Book  `json:"book"`
	Qty  int32 `json:qty`
}

type CartPayment struct {
	Address  string `json:"address"`
	CardNum  string `json:"cardNum"`
	CardCode string `json:"code"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type CartShipment struct {
	Address string `json:"address"`
	Phone   string `json:"phone"`
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
