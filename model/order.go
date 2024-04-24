package model

type OrderStatus string

const (
	OrderStatusInit   OrderStatus = "INIT"
	OrderStatusSold   OrderStatus = "SOLD"
	OrderStatusUnsold OrderStatus = "UNSOLD"
)

// Item provides base information for an item we put to the auction
type Item struct {
	Name         string
	ReservePrice float32
}

// Order provides auction order information
type Order struct {
	Item         Item
	CreationTime int64
	Status       OrderStatus
	CloseTime    int64
	LastBid      float32
	CloseBid     float32
}
