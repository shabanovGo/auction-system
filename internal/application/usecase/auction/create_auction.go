package auction

import (
    "context"
    "auction-system/internal/application/dto/auction"
    "auction-system/internal/domain/repository"
    "auction-system/internal/domain/errors"
)

type CreateAuctionUseCase struct {
    auctionRepo repository.AuctionRepository
    lotRepo     repository.LotRepository
}

func NewCreateAuctionUseCase(
    auctionRepo repository.AuctionRepository,
    lotRepo repository.LotRepository,
) *CreateAuctionUseCase {
    return &CreateAuctionUseCase{
        auctionRepo: auctionRepo,
        lotRepo:     lotRepo,
    }
}

func (uc *CreateAuctionUseCase) Execute(ctx context.Context, req *auction.CreateAuctionRequest) (*auction.AuctionResponse, error) {
    _, err := uc.lotRepo.GetByID(ctx, req.LotID)
    if err != nil {
        return nil, err
    }

    existingAuction, _ := uc.auctionRepo.GetByLotID(ctx, req.LotID)
    if existingAuction != nil {
        return nil, errors.New(errors.ErrorTypeValidation, "auction already exists for this lot", nil)
    }

    auctionEntity := auction.ToEntity(req)
    
    if err := uc.auctionRepo.Create(ctx, auctionEntity); err != nil {
        return nil, err
    }

    return auction.FromEntity(auctionEntity), nil
}
