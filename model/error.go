package model

import (
	"errors"
)

// list of common errors
var (
	ErrUnknownCommandType      = errors.New("unknown command type")
	ErrInvalidData             = errors.New("invalid data")
	ErrNotFound                = errors.New("not found")
	ErrAuctionIsFinishedByTime = errors.New("auction is finished by time")
)
