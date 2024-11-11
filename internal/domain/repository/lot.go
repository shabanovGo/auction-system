package repository

import (
    "context"
    "auction-system/internal/domain/entity"
)

type LotRepository interface {
    Create(ctx context.Context, lot *entity.Lot) error
    GetByID(ctx context.Context, id int64) (*entity.Lot, error)
    Update(ctx context.Context, id int64, lot *entity.Lot) (*entity.Lot, error)
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, offset, limit int) ([]*entity.Lot, int64, error)
}
