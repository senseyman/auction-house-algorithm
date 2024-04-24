package inmemory

import (
	"context"
	"math"
	"sort"
	"sync"

	"github.com/senseyman/auction-house/model"
)

// Storage emulates in-memory storage for storing and processing auction data.
type Storage struct {
	mx sync.Mutex

	orders         map[string]*model.Order         // imitate order table, key - order name
	auctionHistory map[string][]*model.OrderAction // imitate auction_history, key - order name, value - array of auction states
}

func New() *Storage {
	return &Storage{
		orders:         make(map[string]*model.Order),
		auctionHistory: make(map[string][]*model.OrderAction),
	}
}

// CreateOrder method stores order and initiate first auction state
func (s *Storage) CreateOrder(_ context.Context, order model.Order) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.orders[order.Item.Name] = &order

	auctionHistory := s.auctionHistory[order.Item.Name]
	auctionHistory = append(auctionHistory, &model.OrderAction{
		Order:    &order,
		UserID:   0,
		BidValue: 0,
	})
	s.auctionHistory[order.Item.Name] = auctionHistory

	return nil
}

// BidOrder method checks bid value for the order and saves the bid to the history data
func (s *Storage) BidOrder(_ context.Context, bid model.BidCommand) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	// find order
	order, ok := s.orders[bid.ItemName]
	if !ok {
		return model.ErrNotFound
	}

	// check timestamp
	if bid.Timestamp > order.CloseTime {
		return model.ErrAuctionIsFinishedByTime
	}

	// check bid price
	if bid.BidAmount > order.LastBid {
		// update new bid amount
		order.LastBid = bid.BidAmount
	}
	// update order
	s.orders[bid.ItemName] = order

	// update history
	auctionHistory := s.auctionHistory[bid.ItemName]
	auctionHistory = append(auctionHistory, &model.OrderAction{
		Order:    order,
		UserID:   bid.UserID,
		BidValue: bid.BidAmount,
	})
	s.auctionHistory[bid.ItemName] = auctionHistory

	return nil
}

// FinishExpiredAuctions method finishes auctions that expired by time
func (s *Storage) FinishExpiredAuctions(_ context.Context, timestamp int64) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	for orderName := range s.orders {
		order := s.orders[orderName]
		if timestamp > order.CloseTime {
			s.closeOrder(order)
			s.orders[orderName] = order
		}
	}

	return nil
}

// FinishAllAuctions method finishes unfinished auctions
func (s *Storage) FinishAllAuctions(_ context.Context) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	for orderName := range s.orders {
		order := s.orders[orderName]
		if order.Status == model.OrderStatusInit {
			// finish order that is still opened
			s.closeOrder(order)
			s.orders[orderName] = order
		}
	}

	return nil
}

func (s *Storage) closeOrder(order *model.Order) {
	// closing order
	if order.LastBid >= order.Item.ReservePrice {
		order.Status = model.OrderStatusSold
		// set the price
		order.CloseBid = s.getOrderAuctionFinalPrice(order.Item.Name, order.Item.ReservePrice)
	} else {
		order.Status = model.OrderStatusUnsold
	}
}

func (s *Storage) getOrderAuctionFinalPrice(orderName string, reservePrice float32) float32 {
	auctionHistory := s.auctionHistory[orderName]
	if len(auctionHistory) < 3 {
		// we have only initial state and maybe first bid - return reserve price
		return reservePrice
	}

	return auctionHistory[len(auctionHistory)-2].BidValue // return second last bid
}

// GetAuctionResults provides results of all auctions
func (s *Storage) GetAuctionResults(_ context.Context) ([]model.ActionResult, error) {
	results := make([]model.ActionResult, 0, len(s.orders))

	// collect all orders to sort it in a time order
	orders := make([]*model.Order, 0, len(s.orders))
	for orderName := range s.orders {
		orders = append(orders, s.orders[orderName])
	}

	for _, order := range orders {
		userID, stat := s.getOrderAuctionStatistics(order.Item.Name)
		if order.Status == model.OrderStatusUnsold {
			userID = 0
		}
		results = append(results, model.ActionResult{
			CreationTime: order.CreationTime,
			CloseTime:    order.CloseTime,
			Item:         order.Item.Name,
			UserID:       userID,
			Status:       order.Status,
			PricePaid:    order.CloseBid,
			Statistics:   stat,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreationTime < results[j].CreationTime
	})

	return results, nil
}

func (s *Storage) getOrderAuctionStatistics(orderName string) (int, model.AuctionStatistics) {
	auctionHistory := s.auctionHistory[orderName]

	if len(auctionHistory) < 2 { // has only init state
		return 0, model.AuctionStatistics{}
	}

	var (
		maxBid   = float32(0.0)
		minBid   = float32(math.MaxFloat32)
		bidCount = 0
	)
	for idx := 1; idx < len(auctionHistory); idx++ { // skip a first element as an INIT state
		if auctionHistory[idx].BidValue > maxBid {
			maxBid = auctionHistory[idx].BidValue
		}
		if auctionHistory[idx].BidValue < minBid {
			minBid = auctionHistory[idx].BidValue
		}
		bidCount++
	}

	return auctionHistory[len(auctionHistory)-1].UserID, model.AuctionStatistics{
		TotalBidCount: bidCount,
		HighestBid:    maxBid,
		LowestBid:     minBid,
	}
}
