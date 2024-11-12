package repository

import (
    "context"
    "auction-system/internal/domain/entity"
)

type AuctionRepository interface {
    Create(ctx context.Context, auction *entity.Auction) error
    GetByID(ctx context.Context, id int64) (*entity.Auction, error)
    Update(ctx context.Context, id int64, auction *entity.Auction) (*entity.Auction, error)
    List(ctx context.Context, offset, limit int, status *entity.AuctionStatus) ([]*entity.Auction, int64, error)
    
    GetByLotID(ctx context.Context, lotID int64) (*entity.Auction, error)
    UpdateStatus(ctx context.Context, id int64, status entity.AuctionStatus) error
    GetActiveAuctions(ctx context.Context) ([]*entity.Auction, error)
    GetPendingAuctionsToStart(ctx context.Context) ([]*entity.Auction, error)
    GetEndedAuctions(ctx context.Context) ([]*entity.Auction, error)
}
