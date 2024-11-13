package notification

import (
    "context"
    "log"
    "os"
    "auction-system/internal/domain/entity"
    "auction-system/internal/domain/notification"
)

type MockNotificationAdapter struct {
    logger *log.Logger
}

func NewMockNotificationAdapter() notification.NotificationService {
    return &MockNotificationAdapter{
        logger: log.New(os.Stdout, "[NOTIFICATION] ", log.LstdFlags),
    }
}

func (a *MockNotificationAdapter) NotifyAuctionStarted(ctx context.Context, auction *entity.Auction, participants []int64) error {
    a.logger.Printf("Auction %d started. Notifying participants: %v", auction.ID, participants)
    return nil
}

func (a *MockNotificationAdapter) NotifyAuctionResults(ctx context.Context, auction *entity.Auction, winner int64, participants []int64) error {
    a.logger.Printf("Auction %d ended. Winner: %d. Notifying all participants: %v", auction.ID, winner, participants)
    return nil
}

func (a *MockNotificationAdapter) NotifyTransactionStatus(ctx context.Context, userID int64, amount float64, success bool) error {
    status := "succeeded"
    if !success {
        status = "failed"
    }
    a.logger.Printf("Transaction %s for user %d: %.2f", status, userID, amount)
    return nil
}
