package auction

import "auction-system/internal/domain/entity"

func FromEntity(auction *entity.Auction) *AuctionResponse {
    return &AuctionResponse{
        ID:           auction.ID,
        LotID:        auction.LotID,
        StartPrice:   auction.StartPrice,
        MinStep:      auction.MinStep,
        CurrentPrice: auction.CurrentPrice,
        StartTime:    auction.StartTime,
        EndTime:      auction.EndTime,
        Status:       auction.Status,
        WinnerID:     auction.WinnerID,
        WinnerBidID:  auction.WinnerBidID,
        CreatedAt:    auction.CreatedAt,
        UpdatedAt:    auction.UpdatedAt,
    }
}

func ToEntity(req *CreateAuctionRequest) *entity.Auction {
    return &entity.Auction{
        LotID:        req.LotID,
        StartPrice:   req.StartPrice,
        MinStep:      req.MinStep,
        CurrentPrice: req.StartPrice,
        StartTime:    req.StartTime,
        EndTime:      req.EndTime,
        Status:       entity.AuctionStatusPending,
    }
}
