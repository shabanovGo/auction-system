package user

import (
    "context"
    userDto "auction-system/internal/application/dto/user"
    "auction-system/internal/domain/repository"
)

type GetUserUseCase struct {
    userRepo repository.UserRepository
}

func NewGetUserUseCase(userRepo repository.UserRepository) *GetUserUseCase {
    return &GetUserUseCase{
        userRepo: userRepo,
    }
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id int64) (*userDto.UserResponse, error) {
    user, err := uc.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    return userDto.FromEntity(user), nil
}
