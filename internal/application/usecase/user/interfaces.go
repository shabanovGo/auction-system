package user

import (
    "context"
    dto "auction-system/internal/application/dto/user"
)

type CreateUserUseCaseInterface interface {
    Execute(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
}

type GetUserUseCaseInterface interface {
    Execute(ctx context.Context, id int64) (*dto.UserResponse, error)
}

type UpdateUserUseCaseInterface interface {
    Execute(ctx context.Context, id int64, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
}

type DeleteUserUseCaseInterface interface {
    Execute(ctx context.Context, id int64) error
}

type GetAllUsersUseCaseInterface interface {
    Execute(ctx context.Context) (*dto.UserListResponse, error)
}

type UpdateBalanceUseCaseInterface interface {
    Execute(ctx context.Context, input UpdateBalanceInput) error
}
