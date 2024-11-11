package user

import (
    "auction-system/internal/domain/entity"
)

func (r *CreateUserRequest) ToEntity() *entity.User {
    return &entity.User{
        Username: r.Username,
        Email:    r.Email,
        Balance:  r.Balance,
    }
}

func FromEntity(user *entity.User) *UserResponse {
    return &UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Email:     user.Email,
        Balance:   user.Balance,
        CreatedAt: user.CreatedAt,
    }
}

func ToUserListResponse(users []*entity.User, total int64) *UserListResponse {
    response := &UserListResponse{
        Users: make([]UserResponse, len(users)),
        Total: total,
    }

    for i, user := range users {
        response.Users[i] = *FromEntity(user)
    }

    return response
}
