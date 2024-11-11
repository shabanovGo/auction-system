package user

import (
    "context"
    "auction-system/internal/application/dto/user"
	"auction-system/internal/domain/repository"
)

type UpdateUserUseCase struct {
    userRepo repository.UserRepository
}

func NewUpdateUserUseCase(userRepo repository.UserRepository) *UpdateUserUseCase {
    return &UpdateUserUseCase{
        userRepo: userRepo,
    }
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, id int64, req *user.UpdateUserRequest) (*user.UserResponse, error) {
    existingUser, err := uc.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    if req.Username != "" {
        existingUser.Username = req.Username
    }
    if req.Email != "" {
        existingUser.Email = req.Email
    }

    updatedUser, err := uc.userRepo.Update(ctx, id, existingUser)
    if err != nil {
        return nil, err
    }

    return user.FromEntity(updatedUser), nil
}
