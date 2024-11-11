package lot

import (
    "context"
    "auction-system/internal/application/dto/lot"
    "auction-system/internal/domain/repository"
)

type GetLotUseCase struct {
    lotRepo repository.LotRepository
}

func NewGetLotUseCase(lotRepo repository.LotRepository) *GetLotUseCase {
    return &GetLotUseCase{
        lotRepo: lotRepo,
    }
}

func (uc *GetLotUseCase) Execute(ctx context.Context, id int64) (*lot.LotResponse, error) {
    lotEntity, err := uc.lotRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    return lot.FromEntity(lotEntity), nil
}
