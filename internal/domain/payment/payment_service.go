package payment

import (
    "context"
)

type PaymentService interface {
    ProcessPayment(ctx context.Context, userID int64, amount float64) error
    RefundPayment(ctx context.Context, userID int64, amount float64) error
}
