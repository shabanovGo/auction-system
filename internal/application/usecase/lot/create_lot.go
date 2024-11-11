package lot

import (
    "context"
    "auction-system/internal/application/dto/lot"
    "auction-system/internal/domain/repository"
)

type CreateLotUseCase struct {
    lotRepo repository.LotRepository
}

func NewCreateLotUseCase(lotRepo repository.LotRepository) *CreateLotUseCase {
    return &CreateLotUseCase{
        lotRepo: lotRepo,
    }
}

func (uc *CreateLotUseCase) Execute(ctx context.Context, req *lot.CreateLotRequest) (*lot.LotResponse, error) {
    lotEntity := req.ToEntity()
    
    if err := uc.lotRepo.Create(ctx, lotEntity); err != nil {
        return nil, err
    }

    return lot.FromEntity(lotEntity), nil
}
