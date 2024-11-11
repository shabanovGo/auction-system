package user

import (
    "context"
    "auction-system/internal/application/dto/user"
	"auction-system/internal/domain/repository"
	"auction-system/internal/domain/errors"
)

type CreateUserUseCase struct {
    userRepo repository.UserRepository
}

func NewCreateUserUseCase(userRepo repository.UserRepository) *CreateUserUseCase {
    return &CreateUserUseCase{
        userRepo: userRepo,
    }
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error) {
    existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
    if err == nil && existingUser != nil {
        return nil, errors.ErrUserAlreadyExists
    }

    userEntity := req.ToEntity()

    if err := uc.userRepo.Create(ctx, userEntity); err != nil {
        return nil, err
    }

    return user.FromEntity(userEntity), nil
}
