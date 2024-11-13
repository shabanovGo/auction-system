package worker

import (
    "context"
    "log"
    "time"
    "auction-system/internal/domain/entity"
    "auction-system/internal/domain/repository"
    "auction-system/internal/domain/notification"
)

type AuctionStartWorker struct {
    auctionRepo repository.AuctionRepository
    bidRepo     repository.BidRepository
    lotRepo     repository.LotRepository
    notifier    notification.NotificationService
    interval    time.Duration
}

func NewAuctionStartWorker(
    auctionRepo repository.AuctionRepository,
    bidRepo repository.BidRepository,
    lotRepo repository.LotRepository,
    notifier notification.NotificationService,
) *AuctionStartWorker {
    return &AuctionStartWorker{
        auctionRepo: auctionRepo,
        bidRepo:     bidRepo,
        lotRepo:     lotRepo,
        notifier:    notifier,
        interval:    time.Second * 30,
    }
}

func (w *AuctionStartWorker) Start(ctx context.Context) error {
    ticker := time.NewTicker(w.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return nil
        case <-ticker.C:
            if err := w.processAuctions(ctx); err != nil {
                log.Printf("Error processing auctions for start: %v", err)
            }
        }
    }
}

func (w *AuctionStartWorker) processAuctions(ctx context.Context) error {
    auctions, err := w.auctionRepo.GetPendingAuctionsToStart(ctx)
    if err != nil {
        return err
    }

    for _, auction := range auctions {
        lot, err := w.lotRepo.GetByID(ctx, auction.LotID)
        if err != nil {
            log.Printf("Error getting lot for auction %d: %v", auction.ID, err)
            continue
        }

        participants := []int64{lot.CreatorID}

        bidders, err := w.bidRepo.GetUniqueParticipantsByAuctionID(ctx, auction.ID)
        if err != nil {
            log.Printf("Error getting bidders for auction %d: %v", auction.ID, err)
        } else {
            participants = append(participants, bidders...)
        }

        err = w.auctionRepo.UpdateStatus(ctx, auction.ID, entity.AuctionStatusActive)
        if err != nil {
            log.Printf("Error updating auction %d status to ACTIVE: %v", auction.ID, err)
            continue
        }

        if err := w.notifier.NotifyAuctionStarted(ctx, auction, participants); err != nil {
            log.Printf("Error sending notifications for auction %d: %v", auction.ID, err)
        }

        log.Printf("Auction %d started successfully", auction.ID)
    }

    return nil
}
