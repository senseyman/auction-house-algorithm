package auction

import (
	"context"
	"sync"

	"github.com/senseyman/auction-house/model"
)

// Service provides all methods for managing the auctions
type Service struct {
	storage       Storage
	readService   ReadService
	reportService ReportService

	errCh     chan error
	commandCh chan model.Command // channel for forwarding commands from file to processing flow
}

func New(storage Storage, readService ReadService, reportService ReportService) *Service {
	return &Service{
		storage:       storage,
		readService:   readService,
		reportService: reportService,
		errCh:         make(chan error, 10),
		commandCh:     make(chan model.Command),
	}
}

// GetErrChannel returns error channel for processing errors on the app top level
func (s *Service) GetErrChannel() chan error {
	return s.errCh
}

// Start starts the main auction flow.
// It runs reading input file through reader service,
// runs processing all data from the file,
// and reporting results.
func (s *Service) Start(ctx context.Context, filename string) error {
	var wg = &sync.WaitGroup{}

	// run processing commands
	s.run(ctx, wg, s.commandCh)

	// read input file
	if err := s.readService.Read(filename, s.commandCh); err != nil {
		return err
	}

	// wait until all data are processed
	wg.Wait()

	// finish all unfinished auctions in the end
	if err := s.finishAllAuctions(ctx); err != nil {
		return err
	}

	// getting auction results
	finalOrderStatuses, err := s.getAuctionResults(ctx)
	if err != nil {
		return err
	}

	// reporting the results
	return s.reportService.Report(finalOrderStatuses)
}

// run starts thread for processing commands from the file
func (s *Service) run(ctx context.Context, wg *sync.WaitGroup,
	commandCh chan model.Command) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done(): // stop thread by ending global context
				return
			case cmd, ok := <-commandCh:
				if !ok {
					return
				}
				s.processCommand(ctx, cmd)
			}
		}
	}()
}

// processCommand manages command types and calls the appropriate method
func (s *Service) processCommand(ctx context.Context, cmd model.Command) {
	var err error
	switch cmd.Type {
	case model.CommandTypeSell:
		err = s.processSell(ctx, cmd)
	case model.CommandTypeBid:
		err = s.processBid(ctx, cmd)
	case model.CommandTypeHeartbeat:
		err = s.processHeartbeat(ctx, cmd)
	default:
		err = model.ErrUnknownCommandType
	}

	if err != nil {
		s.errCh <- err
	}
}

// processSell processes creating new auction order (someone wants to sell something)
func (s *Service) processSell(ctx context.Context, cmd model.Command) error {
	if cmd.Sell == nil {
		return model.ErrInvalidData
	}
	order := s.newOrder(*cmd.Sell)
	return s.storage.CreateOrder(ctx, order)
}

// processBid processes bids
func (s *Service) processBid(ctx context.Context, cmd model.Command) error {
	if cmd.Bid == nil {
		return model.ErrInvalidData
	}

	return s.storage.BidOrder(ctx, *cmd.Bid)
}

// processHeartbeat processes heartbeat. We check do we need to finish some orders by the time
func (s *Service) processHeartbeat(ctx context.Context, cmd model.Command) error {
	if cmd.Heartbeat == nil {
		return model.ErrInvalidData
	}

	return s.storage.FinishExpiredAuctions(ctx, cmd.Heartbeat.Timestamp)
}

// newOrder makes new order instance
func (s *Service) newOrder(sellOrder model.SellCommand) model.Order {
	return model.Order{
		Item: model.Item{
			Name:         sellOrder.ItemName,
			ReservePrice: sellOrder.ReservePrice,
		},
		CreationTime: sellOrder.Timestamp,
		Status:       model.OrderStatusInit,
		CloseTime:    sellOrder.CloseTime,
	}
}

// finishAllAuctions call the flow to finish all auctions that is still opened
func (s *Service) finishAllAuctions(ctx context.Context) error {
	return s.storage.FinishAllAuctions(ctx)
}

// getAuctionResults returns auction results for the next processing
func (s *Service) getAuctionResults(ctx context.Context) ([]model.ActionResult, error) {
	return s.storage.GetAuctionResults(ctx)
}
