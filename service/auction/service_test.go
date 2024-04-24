package auction

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/senseyman/auction-house/model"
	"github.com/senseyman/auction-house/service/auction/mock"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, New(nil, nil, nil))
}

func TestService_GetErrChannel(t *testing.T) {
	assert.NotNil(t, New(nil, nil, nil).GetErrChannel())
}

func TestService_Start(t *testing.T) {
	filename := "file.txt"
	ar := []model.ActionResult{
		{
			CreationTime: 10,
			CloseTime:    20,
			Item:         "phone_1",
			UserID:       3,
			Status:       "SOLD",
			PricePaid:    12.34,
			Statistics: model.AuctionStatistics{
				TotalBidCount: 3,
				HighestBid:    13.45,
				LowestBid:     9.31,
			},
		},
	}
	testErr := errors.New("test error")

	testCases := []struct {
		name   string
		init   func(t *testing.T, ctx context.Context, filename string) *Service
		hasErr bool
	}{
		{
			name: "success/data/sell",
			init: func(t *testing.T, ctx context.Context, filename string) *Service {
				ctrl := gomock.NewController(t)
				storage := mock.NewMockStorage(ctrl)
				reader := mock.NewMockReadService(ctrl)
				reporter := mock.NewMockReportService(ctrl)

				s := New(storage, reader, reporter)

				sellCmd := model.SellCommand{
					Timestamp:    10,
					UserID:       1,
					ItemName:     "phone_1",
					ReservePrice: 10.10,
					CloseTime:    20,
				}
				order := model.Order{
					Item: model.Item{
						Name:         sellCmd.ItemName,
						ReservePrice: sellCmd.ReservePrice,
					},
					CreationTime: sellCmd.Timestamp,
					Status:       model.OrderStatusInit,
					CloseTime:    sellCmd.CloseTime,
				}

				reader.EXPECT().Read(filename, s.commandCh).Return(nil)
				storage.EXPECT().CreateOrder(ctx, order).Return(nil)
				storage.EXPECT().FinishAllAuctions(ctx).Return(nil)
				storage.EXPECT().GetAuctionResults(ctx).Return(ar, nil)
				reporter.EXPECT().Report(ar).Return(nil)

				go func() {
					// imitate sending parsed sell command
					s.commandCh <- model.Command{
						Type: model.CommandTypeSell,
						Sell: &sellCmd,
					}
					time.Sleep(time.Second * 1)
					close(s.commandCh) // imitate finishing of file reading
				}()

				return s
			},
			hasErr: false,
		},
		{
			name: "success/data/bid",
			init: func(t *testing.T, ctx context.Context, filename string) *Service {
				ctrl := gomock.NewController(t)
				storage := mock.NewMockStorage(ctrl)
				reader := mock.NewMockReadService(ctrl)
				reporter := mock.NewMockReportService(ctrl)

				s := New(storage, reader, reporter)

				bidCmd := model.BidCommand{
					Timestamp: 11,
					UserID:    3,
					ItemName:  "phone_1",
					BidAmount: 10.34,
				}

				reader.EXPECT().Read(filename, s.commandCh).Return(nil)
				storage.EXPECT().BidOrder(ctx, bidCmd).Return(nil)
				storage.EXPECT().FinishAllAuctions(ctx).Return(nil)
				storage.EXPECT().GetAuctionResults(ctx).Return(ar, nil)
				reporter.EXPECT().Report(ar).Return(nil)

				go func() {
					// imitate sending parsed sell command
					s.commandCh <- model.Command{
						Type: model.CommandTypeBid,
						Bid:  &bidCmd,
					}
					time.Sleep(time.Second * 1)
					close(s.commandCh) // imitate finishing of file reading
				}()

				return s
			},
			hasErr: false,
		},
		{
			name: "success/data/heartbeat",
			init: func(t *testing.T, ctx context.Context, filename string) *Service {
				ctrl := gomock.NewController(t)
				storage := mock.NewMockStorage(ctrl)
				reader := mock.NewMockReadService(ctrl)
				reporter := mock.NewMockReportService(ctrl)

				s := New(storage, reader, reporter)

				heartbeatCmd := model.HeartbeatCommand{
					Timestamp: 15,
				}

				reader.EXPECT().Read(filename, s.commandCh).Return(nil)
				storage.EXPECT().FinishExpiredAuctions(ctx, heartbeatCmd.Timestamp).Return(nil)
				storage.EXPECT().FinishAllAuctions(ctx).Return(nil)
				storage.EXPECT().GetAuctionResults(ctx).Return(ar, nil)
				reporter.EXPECT().Report(ar).Return(nil)

				go func() {
					// imitate sending parsed sell command
					s.commandCh <- model.Command{
						Type:      model.CommandTypeHeartbeat,
						Heartbeat: &heartbeatCmd,
					}
					time.Sleep(time.Second * 1)
					close(s.commandCh) // imitate finishing of file reading
				}()

				return s
			},
			hasErr: false,
		},
		{
			name: "success/no_data",
			init: func(t *testing.T, ctx context.Context, filename string) *Service {
				ctrl := gomock.NewController(t)
				storage := mock.NewMockStorage(ctrl)
				reader := mock.NewMockReadService(ctrl)
				reporter := mock.NewMockReportService(ctrl)

				s := New(storage, reader, reporter)

				reader.EXPECT().Read(filename, s.commandCh).Return(nil)
				storage.EXPECT().FinishAllAuctions(ctx).Return(nil)
				storage.EXPECT().GetAuctionResults(ctx).Return(ar, nil)
				reporter.EXPECT().Report(ar).Return(nil)

				go func() {
					time.Sleep(time.Second * 1)
					close(s.commandCh) // imitate finishing of file reading
				}()

				return s
			},
			hasErr: false,
		},
		{
			name: "err/report",
			init: func(t *testing.T, ctx context.Context, filename string) *Service {
				ctrl := gomock.NewController(t)
				storage := mock.NewMockStorage(ctrl)
				reader := mock.NewMockReadService(ctrl)
				reporter := mock.NewMockReportService(ctrl)

				s := New(storage, reader, reporter)

				reader.EXPECT().Read(filename, s.commandCh).Return(nil)
				storage.EXPECT().FinishAllAuctions(ctx).Return(nil)
				storage.EXPECT().GetAuctionResults(ctx).Return(ar, nil)
				reporter.EXPECT().Report(ar).Return(testErr)

				go func() {
					time.Sleep(time.Second * 1)
					close(s.commandCh) // imitate finishing of file reading
				}()

				return s
			},
			hasErr: true,
		},
		{
			name: "err/results",
			init: func(t *testing.T, ctx context.Context, filename string) *Service {
				ctrl := gomock.NewController(t)
				storage := mock.NewMockStorage(ctrl)
				reader := mock.NewMockReadService(ctrl)

				s := New(storage, reader, nil)

				reader.EXPECT().Read(filename, s.commandCh).Return(nil)
				storage.EXPECT().FinishAllAuctions(ctx).Return(nil)
				storage.EXPECT().GetAuctionResults(ctx).Return(nil, testErr)

				go func() {
					time.Sleep(time.Second * 1)
					close(s.commandCh) // imitate finishing of file reading
				}()

				return s
			},
			hasErr: true,
		},
		{
			name: "err/finish_auctions",
			init: func(t *testing.T, ctx context.Context, filename string) *Service {
				ctrl := gomock.NewController(t)
				storage := mock.NewMockStorage(ctrl)
				reader := mock.NewMockReadService(ctrl)

				s := New(storage, reader, nil)

				reader.EXPECT().Read(filename, s.commandCh).Return(nil)
				storage.EXPECT().FinishAllAuctions(ctx).Return(testErr)

				go func() {
					time.Sleep(time.Second * 1)
					close(s.commandCh) // imitate finishing of file reading
				}()

				return s
			},
			hasErr: true,
		},
		{
			name: "err/reader",
			init: func(t *testing.T, ctx context.Context, filename string) *Service {
				ctrl := gomock.NewController(t)
				storage := mock.NewMockStorage(ctrl)
				reader := mock.NewMockReadService(ctrl)

				s := New(storage, reader, nil)

				reader.EXPECT().Read(filename, s.commandCh).Return(testErr)

				go func() {
					time.Sleep(time.Second * 1)
					close(s.commandCh) // imitate finishing of file reading
				}()

				return s
			},
			hasErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			s := tc.init(t, ctx, filename)
			err := s.Start(ctx, filename)

			if tc.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
