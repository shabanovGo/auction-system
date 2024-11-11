package errors

import "errors"

// Domain errors
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user with this email already exists")
    ErrInsufficientBalance = errors.New("insufficient balance")
    ErrInvalidInput     = errors.New("invalid input data")
    ErrLotNotFound = errors.New("lot not found")
    ErrInvalidPage = errors.New("invalid page parameters")
)

// Error types
type ErrorType string

const (
    ErrorTypeNotFound      ErrorType = "NOT_FOUND"
    ErrorTypeValidation    ErrorType = "VALIDATION"
    ErrorTypeConflict      ErrorType = "CONFLICT"
    ErrorTypeInternal      ErrorType = "INTERNAL"
    ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"
)

// AppError представляет ошибку приложения
type AppError struct {
    Type    ErrorType `json:"type"`
    Message string    `json:"message"`
    Err     error    `json:"-"`
}

func (e *AppError) Error() string {
    return e.Message
}

// New создает новую ошибку приложения
func New(errType ErrorType, message string, err error) *AppError {
    return &AppError{
        Type:    errType,
        Message: message,
        Err:     err,
    }
}

// NewValidationError создает ошибку валидации
func NewValidationError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeValidation,
        Message: message,
    }
}

// NewNotFoundError создает ошибку "не найдено"
func NewNotFoundError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeNotFound,
        Message: message,
    }
}

// NewConflictError создает ошибку конфликта
func NewConflictError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeConflict,
        Message: message,
    }
}
