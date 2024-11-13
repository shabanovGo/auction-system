package bid

import "time"

type BidResponse struct {
    ID        int64     `json:"id"`
    AuctionID int64     `json:"auction_id"`
    UserID    int64     `json:"user_id"`
    Amount    float64   `json:"amount"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ListBidsResponse struct {
    Bids       []BidResponse `json:"bids"`
    TotalCount int64         `json:"total_count"`
}
