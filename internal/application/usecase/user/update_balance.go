package user

import (
    "context"
    "auction-system/internal/domain/repository"
	"auction-system/internal/domain/errors"
)

type UpdateBalanceUseCase struct {
    userRepo repository.UserRepository
}

func NewUpdateBalanceUseCase(userRepo repository.UserRepository) *UpdateBalanceUseCase {
    return &UpdateBalanceUseCase{
        userRepo: userRepo,
    }
}

type UpdateBalanceInput struct {
    UserID  int64   `json:"user_id"`
    Amount  float64 `json:"amount"`
}

func (uc *UpdateBalanceUseCase) Execute(ctx context.Context, input UpdateBalanceInput) error {
    user, err := uc.userRepo.GetByID(ctx, input.UserID)
    if err != nil {
        return err
    }

    newBalance := user.Balance + input.Amount
    if newBalance < 0 {
        return errors.ErrInsufficientBalance
    }

    return uc.userRepo.UpdateBalance(ctx, input.UserID, newBalance)
}
