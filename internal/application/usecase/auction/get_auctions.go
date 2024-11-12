package auction

import (
    "context"
    "auction-system/internal/application/dto/auction"
    "auction-system/internal/domain/repository"
    "auction-system/internal/domain/entity"
)

type ListAuctionsUseCase struct {
    auctionRepo repository.AuctionRepository
}

func NewListAuctionsUseCase(auctionRepo repository.AuctionRepository) *ListAuctionsUseCase {
    return &ListAuctionsUseCase{
        auctionRepo: auctionRepo,
    }
}

func (uc *ListAuctionsUseCase) Execute(ctx context.Context, page, pageSize int, status *string) (*auction.ListAuctionsResponse, error) {
    var auctionStatus *entity.AuctionStatus
    if status != nil {
        s := entity.AuctionStatus(*status)
        auctionStatus = &s
    }

    offset := (page - 1) * pageSize
    auctions, total, err := uc.auctionRepo.List(ctx, offset, pageSize, auctionStatus)
    if err != nil {
        return nil, err
    }

    response := &auction.ListAuctionsResponse{
        Auctions:   make([]auction.AuctionResponse, len(auctions)),
        TotalCount: total,
    }

    for i, a := range auctions {
        response.Auctions[i] = *auction.FromEntity(a)
    }

    return response, nil
}
