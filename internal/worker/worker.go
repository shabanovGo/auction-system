package worker

import (
	"context"
	"time"
	"auction-system/internal/domain/repository"
	"auction-system/internal/domain/notification"
)

type Worker struct {
	auctionStartWorker *AuctionStartWorker
	auctionEndWorker   *AuctionCloserWorker
}

func NewWorker(
	auctionRepo repository.AuctionRepository,
	lotRepo repository.LotRepository,
	userRepo repository.UserRepository,
	bidRepo repository.BidRepository,
	notifier notification.NotificationService,
) *Worker {
	return &Worker{
		auctionStartWorker: NewAuctionStartWorker(auctionRepo, bidRepo, lotRepo, notifier),
		auctionEndWorker:   NewAuctionCloserWorker(auctionRepo, bidRepo, userRepo, lotRepo, time.Second * 30),
	}
}

func (w *Worker) Start(ctx context.Context) {
	go w.auctionStartWorker.Start(ctx)
	
	go w.auctionEndWorker.Start(ctx)
}
