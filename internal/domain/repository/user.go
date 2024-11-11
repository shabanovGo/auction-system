package repository

import (
    "context"
    "auction-system/internal/domain/entity"
)

type UserRepository interface {
    GetAll(ctx context.Context) ([]*entity.User, error)
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
    UpdateBalance(ctx context.Context, id int64, balance float64) error
    List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)
	Update(ctx context.Context, id int64, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id int64) error
}
