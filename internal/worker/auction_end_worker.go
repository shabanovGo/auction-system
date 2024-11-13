package worker

import (
    "context"
    "log"
    "time"
    
    "auction-system/internal/domain/entity"
    "auction-system/internal/domain/repository"
)

type AuctionCloserWorker struct {
    auctionRepo repository.AuctionRepository
    bidRepo     repository.BidRepository
    userRepo    repository.UserRepository
    lotRepo     repository.LotRepository
    interval    time.Duration
}

func NewAuctionCloserWorker(
    auctionRepo repository.AuctionRepository,
    bidRepo repository.BidRepository,
    userRepo repository.UserRepository,
    lotRepo repository.LotRepository,
    interval time.Duration,
) *AuctionCloserWorker {
    return &AuctionCloserWorker{
        auctionRepo: auctionRepo,
        bidRepo:     bidRepo,
        userRepo:    userRepo,
        lotRepo:     lotRepo,
        interval:    interval,
    }
}

func (w *AuctionCloserWorker) Start(ctx context.Context) {
    ticker := time.NewTicker(w.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := w.processEndedAuctions(ctx); err != nil {
                log.Printf("Error processing ended auctions: %v", err)
            }
        }
    }
}

func (w *AuctionCloserWorker) processEndedAuctions(ctx context.Context) error {
    endedAuctions, err := w.auctionRepo.GetEndedAuctions(ctx)
    if err != nil {
        return err
    }

    for _, auction := range endedAuctions {
        if err := w.processAuction(ctx, auction); err != nil {
            log.Printf("Error processing auction %d: %v", auction.ID, err)
            continue
        }
    }

    return nil
}

func (w *AuctionCloserWorker) processAuction(ctx context.Context, auction *entity.Auction) error {
    bids, err := w.bidRepo.GetByAuctionID(ctx, auction.ID)
    if err != nil {
        return err
    }

    if len(bids) > 0 {
        winningBid := bids[0]
        
        lot, err := w.lotRepo.GetByID(ctx, auction.LotID)
        if err != nil {
            return err
        }

        winner, err := w.userRepo.GetByID(ctx, winningBid.UserID)
        if err != nil {
            return err
        }

        newWinnerBalance := winner.Balance - winningBid.Amount
        if err := w.userRepo.UpdateBalance(ctx, winner.ID, newWinnerBalance); err != nil {
            return err
        }

        seller, err := w.userRepo.GetByID(ctx, lot.CreatorID)
        if err != nil {
            return err
        }
        
        newSellerBalance := seller.Balance + winningBid.Amount
        if err := w.userRepo.UpdateBalance(ctx, seller.ID, newSellerBalance); err != nil {
            _ = w.userRepo.UpdateBalance(ctx, winner.ID, winner.Balance)
            return err
        }

        auction.Status = entity.AuctionStatusEnded
        auction.WinnerID = &winningBid.UserID
        auction.WinnerBidID = &winningBid.ID
    } else {
        auction.Status = entity.AuctionStatusEnded
    }

    if err := w.auctionRepo.UpdateStatus(ctx, auction.ID, auction.Status); err != nil {
        return err
    }

    return nil
}
