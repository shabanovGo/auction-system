package notification

import (
    "context"
    "auction-system/internal/domain/entity"
)

type NotificationService interface {
    NotifyAuctionStarted(ctx context.Context, auction *entity.Auction, participants []int64) error
    NotifyAuctionResults(ctx context.Context, auction *entity.Auction, winner int64, participants []int64) error
    NotifyTransactionStatus(ctx context.Context, userID int64, amount float64, success bool) error
}
