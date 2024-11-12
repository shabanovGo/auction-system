package auction

import (
    "context"
    "auction-system/internal/application/dto/auction"
    "auction-system/internal/domain/repository"
)

type GetAuctionUseCase struct {
    auctionRepo repository.AuctionRepository
}

func NewGetAuctionUseCase(auctionRepo repository.AuctionRepository) *GetAuctionUseCase {
    return &GetAuctionUseCase{
        auctionRepo: auctionRepo,
    }
}

func (uc *GetAuctionUseCase) Execute(ctx context.Context, id int64) (*auction.AuctionResponse, error) {
    auctionEntity, err := uc.auctionRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    return auction.FromEntity(auctionEntity), nil
}
