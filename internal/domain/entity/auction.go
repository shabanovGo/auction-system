package entity

import (
    "time"
)

type Auction struct {
    ID           int64     `json:"id"`
    LotID        int64     `json:"lot_id"`
    StartPrice   float64   `json:"start_price"`
    MinStep      float64   `json:"min_step"`
    CurrentPrice float64   `json:"current_price"`
    StartTime    time.Time `json:"start_time"`
    EndTime      time.Time `json:"end_time"`
    Status       AuctionStatus `json:"status"`
    WinnerID     *int64    `json:"winner_id,omitempty"`
    WinnerBidID  *int64    `json:"winner_bid_id,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type AuctionStatus string

const (
    AuctionStatusPending   AuctionStatus = "PENDING"
    AuctionStatusActive    AuctionStatus = "ACTIVE"
    AuctionStatusEnded     AuctionStatus = "ENDED"
    AuctionStatusCompleted AuctionStatus = "COMPLETED"
    AuctionStatusCanceled  AuctionStatus = "CANCELED"
)
