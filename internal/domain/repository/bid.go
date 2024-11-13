package repository

import (
    "context"
    "auction-system/internal/domain/entity"
)

type BidRepository interface {
    Create(ctx context.Context, bid *entity.Bid) error
    GetByID(ctx context.Context, id int64) (*entity.Bid, error)
    GetByAuctionID(ctx context.Context, auctionID int64) ([]*entity.Bid, error)
    List(ctx context.Context, auctionID int64, offset, limit int) ([]*entity.Bid, int64, error)
    GetUniqueParticipantsByAuctionID(ctx context.Context, auctionID int64) ([]int64, error)
}
