package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/senseyman/auction-house/service/auction"
	"github.com/senseyman/auction-house/service/reader"
	"github.com/senseyman/auction-house/service/report"
	"github.com/senseyman/auction-house/storage/inmemory"
)

var (
	filePathFlag = flag.String("path", "input.txt", "")
)

func main() {
	flag.Parse()

	// first - init all services
	storage := inmemory.New()
	readService := reader.New()
	reportService := report.New()
	auctionService := auction.New(storage, readService, reportService)

	// create global context with cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupGracefulShutdown(cancel)

	// run thread for processing err messages
	go processErrMsgs(auctionService.GetErrChannel())

	// run the main flow
	if err := auctionService.Start(ctx, *filePathFlag); err != nil {
		fmt.Printf("error while executing auction: %v\n", err)
	}
}

func processErrMsgs(errCh chan error) {
	for range errCh {
		// add printing error if needed
	}
}

// setupGracefulShutdown provides processing income os signals and stopping the app by canceling global app context
func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		fmt.Printf("Got Interrupt signal: %v\n", sig.String())
		stop()
	}()
}
