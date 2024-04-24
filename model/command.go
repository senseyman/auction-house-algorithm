package model

type CommandType int

const (
	CommandTypeUnknown CommandType = iota
	CommandTypeSell
	CommandTypeBid
	CommandTypeHeartbeat
)

// Command struct contains commands from input file for future processing
type Command struct {
	Type      CommandType
	Sell      *SellCommand
	Bid       *BidCommand
	Heartbeat *HeartbeatCommand
}

// SellCommand provides sell instructions
type SellCommand struct {
	Timestamp    int64
	UserID       int
	ItemName     string
	ReservePrice float32
	CloseTime    int64
}

// BidCommand provides bid instructions
type BidCommand struct {
	Timestamp int64
	UserID    int
	ItemName  string
	BidAmount float32
}

// HeartbeatCommand provides heartbeat instructions
type HeartbeatCommand struct {
	Timestamp int64
}
