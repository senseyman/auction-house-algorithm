package auction

import (
	"context"

	"github.com/senseyman/auction-house/model"
)

type Storage interface {
	CreateOrder(ctx context.Context, order model.Order) error
	BidOrder(ctx context.Context, bid model.BidCommand) error
	FinishExpiredAuctions(ctx context.Context, timestamp int64) error
	FinishAllAuctions(_ context.Context) error
	GetAuctionResults(ctx context.Context) ([]model.ActionResult, error)
}

type ReadService interface {
	Read(filename string, outputCh chan model.Command) error
}

type ReportService interface {
	Report(fos []model.ActionResult) error
}
