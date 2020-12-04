package dto

type CartItem struct {
	Book Book  `json:"book"`
	Qty  int32 `json:qty`
}
