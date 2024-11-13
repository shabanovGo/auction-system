package notification

import "time"

type NotificationType string

const (
    AuctionStarted NotificationType = "AUCTION_STARTED"
    AuctionClosed   NotificationType = "AUCTION_CLOSED"
    AuctionWon      NotificationType = "AUCTION_WON"
    NewBid          NotificationType = "NEW_BID"
    TransactionComplete NotificationType = "TRANSACTION_COMPLETE"
    TransactionFailed  NotificationType = "TRANSACTION_FAILED"
)

type Notification struct {
    UserID    int64           `json:"user_id"`
    Type      NotificationType `json:"type"`
    Message   string          `json:"message"`
    CreatedAt time.Time       `json:"created_at"`
}
