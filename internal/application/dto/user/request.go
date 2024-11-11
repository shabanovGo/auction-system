package user

type CreateUserRequest struct {
    Username string  `json:"username" validate:"required,min=3,max=50"`
    Email    string  `json:"email" validate:"required,email"`
    Balance  float64 `json:"balance" validate:"gte=0"`
}

type UpdateBalanceRequest struct {
    Amount float64 `json:"amount" validate:"required"`
}

type UpdateUserRequest struct {
    Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
    Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

type UpdateBalanceInput struct {
    UserID int64   `json:"user_id"`
    Amount float64 `json:"amount"`
}
