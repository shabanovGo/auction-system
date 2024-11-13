package payment

import (
    "context"
    "log"
    "os"
    "auction-system/internal/domain/payment"
)

type MockPaymentAdapter struct {
    logger *log.Logger
}

func NewMockPaymentAdapter() payment.PaymentService {
    return &MockPaymentAdapter{
        logger: log.New(os.Stdout, "[PAYMENT] ", log.LstdFlags),
    }
}

func (a *MockPaymentAdapter) ProcessPayment(ctx context.Context, userID int64, amount float64) error {
    a.logger.Printf("Processing payment: UserID=%d, Amount=%.2f", userID, amount)
    return nil
}

func (a *MockPaymentAdapter) RefundPayment(ctx context.Context, userID int64, amount float64) error {
    a.logger.Printf("Refunding payment: UserID=%d, Amount=%.2f", userID, amount)
    return nil
}
