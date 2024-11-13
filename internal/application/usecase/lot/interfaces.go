package lot

import (
    "context"
    dto "auction-system/internal/application/dto/lot"
)

type CreateLotUseCaseInterface interface {
    Execute(ctx context.Context, req *dto.CreateLotRequest) (*dto.LotResponse, error)
}

type GetLotUseCaseInterface interface {
    Execute(ctx context.Context, id int64) (*dto.LotResponse, error)
}

type UpdateLotUseCaseInterface interface {
    Execute(ctx context.Context, id int64, req *dto.UpdateLotRequest) (*dto.LotResponse, error)
}

type DeleteLotUseCaseInterface interface {
    Execute(ctx context.Context, id int64) error
}

type ListLotsUseCaseInterface interface {
    Execute(ctx context.Context, page, pageSize int) (*dto.LotListResponse, error)
}
