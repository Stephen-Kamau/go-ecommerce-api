package orders

type OrderItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type CreateOrderRequest struct {
	CustomerRef string             `json:"customer_ref"`
	Items       []OrderItemRequest `json:"items"`
}
