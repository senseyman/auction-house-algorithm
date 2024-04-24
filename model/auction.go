package model

// OrderAction provides information about bids by users
type OrderAction struct {
	Order    *Order
	UserID   int
	BidValue float32
}

// ActionResult - auction result for an order
type ActionResult struct {
	CreationTime int64
	CloseTime    int64
	Item         string
	UserID       int
	Status       OrderStatus
	PricePaid    float32
	Statistics   AuctionStatistics
}

// AuctionStatistics provides some helpful statistics about an order auction
type AuctionStatistics struct {
	TotalBidCount int
	HighestBid    float32
	LowestBid     float32
}
