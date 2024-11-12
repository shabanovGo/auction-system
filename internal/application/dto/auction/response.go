package auction

import (
    "time"
    "auction-system/internal/domain/entity"
)

type AuctionResponse struct {
    ID           int64             `json:"id"`
    LotID        int64             `json:"lot_id"`
    StartPrice   float64           `json:"start_price"`
    MinStep      float64           `json:"min_step"`
    CurrentPrice float64           `json:"current_price"`
    StartTime    time.Time         `json:"start_time"`
    EndTime      time.Time         `json:"end_time"`
    Status       entity.AuctionStatus `json:"status"`
    WinnerID     *int64            `json:"winner_id,omitempty"`
    WinnerBidID  *int64            `json:"winner_bid_id,omitempty"`
    CreatedAt    time.Time         `json:"created_at"`
    UpdatedAt    time.Time         `json:"updated_at"`
}

type BidResponse struct {
    ID        int64
    AuctionID int64
    UserID    int64
    Amount    float64
    CreatedAt time.Time
    UpdatedAt time.Time
}

type ListAuctionsResponse struct {
    Auctions   []AuctionResponse `json:"auctions"`
    TotalCount int64            `json:"total_count"`
}
