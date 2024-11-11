package lot

import (
    "context"
    "auction-system/internal/application/dto/lot"
    "auction-system/internal/domain/repository"
)

type UpdateLotUseCase struct {
    lotRepo repository.LotRepository
}

func NewUpdateLotUseCase(lotRepo repository.LotRepository) *UpdateLotUseCase {
    return &UpdateLotUseCase{
        lotRepo: lotRepo,
    }
}

func (uc *UpdateLotUseCase) Execute(ctx context.Context, id int64, req *lot.UpdateLotRequest) (*lot.LotResponse, error) {
    existingLot, err := uc.lotRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    if req.Title != "" {
        existingLot.Title = req.Title
    }
    if req.Description != "" {
        existingLot.Description = req.Description
    }
    if req.StartPrice > 0 {
        existingLot.StartPrice = req.StartPrice
    }

    updatedLot, err := uc.lotRepo.Update(ctx, id, existingLot)
    if err != nil {
        return nil, err
    }

    return lot.FromEntity(updatedLot), nil
}
