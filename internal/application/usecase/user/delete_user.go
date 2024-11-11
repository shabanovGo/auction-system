package user

import (
    "context"
	"auction-system/internal/domain/repository"
)

type DeleteUserUseCase struct {
    userRepo repository.UserRepository
}

func NewDeleteUserUseCase(userRepo repository.UserRepository) *DeleteUserUseCase {
    return &DeleteUserUseCase{
        userRepo: userRepo,
    }
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id int64) error {
    if err := uc.userRepo.Delete(ctx, id); err != nil {
        return err
    }
    return nil
}
