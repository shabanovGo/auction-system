package bid

import (
    "auction-system/internal/domain/entity"
)

func ToEntity(req *PlaceBidRequest) *entity.Bid {
    return &entity.Bid{
        AuctionID: req.AuctionID,
        UserID:    req.UserID,
        Amount:    req.Amount,
    }
}

func FromEntity(bid *entity.Bid) *BidResponse {
    if bid == nil {
        return nil
    }
    return &BidResponse{
        ID:        bid.ID,
        AuctionID: bid.AuctionID,
        UserID:    bid.UserID,
        Amount:    bid.Amount,
        CreatedAt: bid.CreatedAt,
        UpdatedAt: bid.UpdatedAt,
    }
}

func FromEntityList(bids []*entity.Bid) []BidResponse {
    result := make([]BidResponse, len(bids))
    for i, bid := range bids {
        result[i] = *FromEntity(bid)
    }
    return result
}
