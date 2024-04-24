package reader

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/senseyman/auction-house/model"
)

// Service provides function for reading data from the file where we have auction instructions.
type Service struct {
}

func New() *Service {
	return &Service{}
}

// Read reads data from the file and sends commands to the channel
func (s *Service) Read(filename string, outputCh chan model.Command) error {
	defer close(outputCh)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		cmd := s.parseLineToCommand(line)
		outputCh <- cmd
	}

	return nil
}

// parseLineToCommand parses and determines what kind of command do we have
func (s *Service) parseLineToCommand(line string) model.Command {
	elements := strings.Split(line, "|")
	switch len(elements) {
	case 6: // sell command
		return model.Command{
			Type: model.CommandTypeSell,
			Sell: toSellCommand(elements),
		}
	case 5: // bid command
		return model.Command{
			Type: model.CommandTypeBid,
			Bid:  toBidCommand(elements),
		}
	case 1: // heartbeat command
		return model.Command{
			Type:      model.CommandTypeHeartbeat,
			Heartbeat: toHeartbeatCommand(elements),
		}
	default:
		return model.Command{
			Type: model.CommandTypeUnknown,
		}
	}
}

func toSellCommand(elements []string) *model.SellCommand {
	// skipping element index 2 - action. Always SELL
	var (
		timestamp    int64
		userID       int64
		itemName     string
		reservePrice float64
		closeTime    int64
	)
	timestamp, err := strconv.ParseInt(elements[0], 10, 64)
	if err != nil {
		return nil
	}
	userID, err = strconv.ParseInt(elements[1], 10, 64)
	if err != nil {
		return nil
	}
	itemName = elements[3]
	reservePrice, err = strconv.ParseFloat(elements[4], 32)
	if err != nil {
		return nil
	}
	closeTime, err = strconv.ParseInt(elements[5], 10, 64)
	if err != nil {
		return nil
	}

	return &model.SellCommand{
		Timestamp:    timestamp,
		UserID:       int(userID),
		ItemName:     itemName,
		ReservePrice: float32(reservePrice),
		CloseTime:    closeTime,
	}
}

func toBidCommand(elements []string) *model.BidCommand {
	// skipping element index 2 - action. Always SELL
	var (
		timestamp int64
		userID    int64
		itemName  string
		bidAmount float64
	)
	timestamp, err := strconv.ParseInt(elements[0], 10, 64)
	if err != nil {
		return nil
	}
	userID, err = strconv.ParseInt(elements[1], 10, 64)
	if err != nil {
		return nil
	}
	itemName = elements[3]
	bidAmount, err = strconv.ParseFloat(elements[4], 32)
	if err != nil {
		return nil
	}

	return &model.BidCommand{
		Timestamp: timestamp,
		UserID:    int(userID),
		ItemName:  itemName,
		BidAmount: float32(bidAmount),
	}
}

func toHeartbeatCommand(elements []string) *model.HeartbeatCommand {
	timestamp, err := strconv.ParseInt(elements[0], 10, 64)
	if err != nil {
		return nil
	}

	return &model.HeartbeatCommand{Timestamp: timestamp}
}
