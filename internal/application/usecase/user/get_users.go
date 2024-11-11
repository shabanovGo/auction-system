package user

import (
    "context"
    userDto "auction-system/internal/application/dto/user"
    "auction-system/internal/domain/repository"
)

type GetAllUserUseCase struct {
    userRepo repository.UserRepository
}

func NewGetAllUserUseCase(userRepo repository.UserRepository) *GetAllUserUseCase {
    return &GetAllUserUseCase{
        userRepo: userRepo,
    }
}

func (uc *GetAllUserUseCase) Execute(ctx context.Context) (*userDto.UserListResponse, error) {
    users, err := uc.userRepo.GetAll(ctx)
    if err != nil {
        return nil, err
    }

    return userDto.ToUserListResponse(users, int64(len(users))), nil
}
