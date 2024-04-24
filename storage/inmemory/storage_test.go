package inmemory

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/senseyman/auction-house/model"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, New())
}

func TestStorage_CreateOrder(t *testing.T) {
	order := model.Order{
		Item: model.Item{
			Name:         "phone_1",
			ReservePrice: 20.23,
		},
		CreationTime: 10,
		Status:       model.OrderStatusInit,
		CloseTime:    20,
	}
	auctionInitState := model.OrderAction{
		Order:    &order,
		UserID:   0,
		BidValue: 0,
	}

	storage := New()
	err := storage.CreateOrder(context.TODO(), order)
	assert.NoError(t, err)

	assert.Len(t, storage.orders, 1)
	assert.Len(t, storage.auctionHistory, 1)
	assert.EqualValues(t, auctionInitState, *storage.auctionHistory[order.Item.Name][0])
}

func TestStorage_BidOrder(t *testing.T) {
	itemName := "phone_1"
	order := model.Order{
		Item: model.Item{
			Name:         itemName,
			ReservePrice: 20,
		},
		CreationTime: 10,
		Status:       model.OrderStatusInit,
		CloseTime:    20,
	}

	testCases := []struct {
		name     string
		init     func() *Storage
		bidValue model.BidCommand
		hasErr   bool
	}{
		{
			name: "success",
			init: func() *Storage {
				s := New()

				s.CreateOrder(context.TODO(), order)

				return s
			},
			bidValue: model.BidCommand{
				Timestamp: 12,
				UserID:    3,
				ItemName:  itemName,
				BidAmount: 15.45,
			},
			hasErr: false,
		},
		{
			name: "err/to_late_bid",
			init: func() *Storage {
				s := New()

				s.CreateOrder(context.TODO(), order)

				return s
			},
			bidValue: model.BidCommand{
				Timestamp: 23,
				UserID:    3,
				ItemName:  itemName,
				BidAmount: 15.45,
			},
			hasErr: true,
		},
		{
			name: "err/not_found",
			init: func() *Storage {
				s := New()

				return s
			},
			bidValue: model.BidCommand{
				Timestamp: 23,
				UserID:    3,
				ItemName:  itemName,
				BidAmount: 15.45,
			},
			hasErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.init()
			err := s.BidOrder(context.TODO(), tc.bidValue)

			if tc.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				newOrder := order
				newOrder.LastBid = tc.bidValue.BidAmount
				orderAction := model.OrderAction{
					Order:    &newOrder,
					UserID:   tc.bidValue.UserID,
					BidValue: tc.bidValue.BidAmount,
				}
				assert.EqualValues(t, newOrder, *s.orders[order.Item.Name])
				assert.Len(t, s.auctionHistory[order.Item.Name], 2)
				assert.EqualValues(t, orderAction, *s.auctionHistory[order.Item.Name][1])
			}
		})
	}
}

func TestStorage_FinishExpiredAuctions(t *testing.T) {
	orders := generateOrders(3)
	orders[2].CloseTime = orders[0].CloseTime // to make order 3 sold by time

	s := New()
	// create 3 orders
	err := s.CreateOrder(context.TODO(), orders[0])
	assert.NoError(t, err)
	err = s.CreateOrder(context.TODO(), orders[1])
	assert.NoError(t, err)
	err = s.CreateOrder(context.TODO(), orders[2])
	assert.NoError(t, err)

	// check we have 3 orders inside
	assert.Len(t, s.orders, 3)

	// imitate one bid for order 3
	err = s.BidOrder(context.TODO(), model.BidCommand{
		Timestamp: 15,
		UserID:    3,
		ItemName:  "phone_3",
		BidAmount: 22,
	})
	assert.NoError(t, err)

	err = s.FinishExpiredAuctions(context.TODO(), 19)
	assert.NoError(t, err)
	assert.Equal(t, model.OrderStatusUnsold, s.orders[orders[0].Item.Name].Status)
	assert.Equal(t, model.OrderStatusInit, s.orders[orders[1].Item.Name].Status)
	assert.Equal(t, model.OrderStatusSold, s.orders[orders[2].Item.Name].Status)
}

func TestStorage_FinishAllAuctions(t *testing.T) {
	orders := generateOrders(3)
	orders[2].CloseTime = orders[0].CloseTime // to make order 3 sold by time

	s := New()
	// create 3 orders
	err := s.CreateOrder(context.TODO(), orders[0])
	assert.NoError(t, err)
	err = s.CreateOrder(context.TODO(), orders[1])
	assert.NoError(t, err)
	err = s.CreateOrder(context.TODO(), orders[2])
	assert.NoError(t, err)

	// check we have 3 orders inside
	assert.Len(t, s.orders, 3)

	// imitate one bid for order 3
	err = s.BidOrder(context.TODO(), model.BidCommand{
		Timestamp: 15,
		UserID:    3,
		ItemName:  "phone_3",
		BidAmount: 22,
	})
	assert.NoError(t, err)

	err = s.FinishAllAuctions(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, model.OrderStatusUnsold, s.orders[orders[0].Item.Name].Status)
	assert.Equal(t, model.OrderStatusUnsold, s.orders[orders[1].Item.Name].Status)
	assert.Equal(t, model.OrderStatusSold, s.orders[orders[2].Item.Name].Status)
}

func TestStorage_GetAuctionResults(t *testing.T) {
	expRes := []model.ActionResult{
		{
			CreationTime: 10,
			CloseTime:    15,
			Item:         "phone_1",
			UserID:       0,
			Status:       model.OrderStatusUnsold,
			PricePaid:    0,
			Statistics: model.AuctionStatistics{
				TotalBidCount: 0,
				HighestBid:    0,
				LowestBid:     0,
			},
		},
		{
			CreationTime: 11,
			CloseTime:    20,
			Item:         "phone_2",
			UserID:       0,
			Status:       model.OrderStatusUnsold,
			PricePaid:    0,
			Statistics: model.AuctionStatistics{
				TotalBidCount: 0,
				HighestBid:    0,
				LowestBid:     0,
			},
		},
		{
			CreationTime: 12,
			CloseTime:    15,
			Item:         "phone_3",
			UserID:       3,
			Status:       model.OrderStatusSold,
			PricePaid:    20,
			Statistics: model.AuctionStatistics{
				TotalBidCount: 1,
				HighestBid:    22,
				LowestBid:     22,
			},
		},
	}

	orders := generateOrders(3)
	orders[2].CloseTime = orders[0].CloseTime // to make order 3 sold by time

	s := New()
	// create 3 orders
	err := s.CreateOrder(context.TODO(), orders[0])
	assert.NoError(t, err)
	err = s.CreateOrder(context.TODO(), orders[1])
	assert.NoError(t, err)
	err = s.CreateOrder(context.TODO(), orders[2])
	assert.NoError(t, err)

	// check we have 3 orders inside
	assert.Len(t, s.orders, 3)

	// imitate one bid for order 3
	err = s.BidOrder(context.TODO(), model.BidCommand{
		Timestamp: 15,
		UserID:    3,
		ItemName:  "phone_3",
		BidAmount: 22,
	})
	assert.NoError(t, err)

	err = s.FinishAllAuctions(context.TODO())
	assert.NoError(t, err)

	results, err := s.GetAuctionResults(context.TODO())
	assert.NoError(t, err)
	assert.EqualValues(t, expRes, results)
}

func generateOrders(num int) []model.Order {
	res := make([]model.Order, num)
	for idx := range res {
		res[idx] = model.Order{
			Item: model.Item{
				Name:         fmt.Sprintf("phone_%d", idx+1),
				ReservePrice: 20,
			},
			CreationTime: int64(10 + idx),
			Status:       model.OrderStatusInit,
			CloseTime:    int64(10 + ((idx + 1) * 5)),
		}
	}

	return res
}
