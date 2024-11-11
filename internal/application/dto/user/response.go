package user

import "time"

type UserResponse struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Balance   float64   `json:"balance"`
    CreatedAt time.Time `json:"created_at"`
}

type UserListResponse struct {
    Users []UserResponse `json:"users"`
    Total int64          `json:"total"`
}

type BalanceResponse struct {
    UserID    int64     `json:"user_id"`
    Balance   float64   `json:"balance"`
    UpdatedAt time.Time `json:"updated_at"`
}
