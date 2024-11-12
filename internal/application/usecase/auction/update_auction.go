package auction

import (
    "context"
    "auction-system/internal/application/dto/auction"
    "auction-system/internal/domain/repository"
    "auction-system/internal/domain/errors"
    "auction-system/internal/domain/entity"
)

type UpdateAuctionUseCase struct {
    auctionRepo repository.AuctionRepository
}

func NewUpdateAuctionUseCase(auctionRepo repository.AuctionRepository) *UpdateAuctionUseCase {
    return &UpdateAuctionUseCase{
        auctionRepo: auctionRepo,
    }
}

func (uc *UpdateAuctionUseCase) Execute(ctx context.Context, id int64, req *auction.UpdateAuctionRequest) (*auction.AuctionResponse, error) {
    existingAuction, err := uc.auctionRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    if existingAuction.Status != entity.AuctionStatusPending {
        return nil, errors.New(errors.ErrorTypeValidation, "can only update pending auctions", nil)
    }

    updateAuction := &entity.Auction{
        ID: id,
    }

    if req.StartPrice != nil {
        updateAuction.StartPrice = *req.StartPrice
        updateAuction.CurrentPrice = *req.StartPrice
    }
    if req.MinStep != nil {
        updateAuction.MinStep = *req.MinStep
    }
    if req.StartTime != nil {
        updateAuction.StartTime = *req.StartTime
    }
    if req.EndTime != nil {
        updateAuction.EndTime = *req.EndTime
    }
    if req.Status != nil {
        updateAuction.Status = entity.AuctionStatus(*req.Status)
    }

    updatedAuction, err := uc.auctionRepo.Update(ctx, id, updateAuction)
    if err != nil {
        return nil, err
    }

    return auction.FromEntity(updatedAuction), nil
}
