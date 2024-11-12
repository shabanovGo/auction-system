package auction

import "time"

type CreateAuctionRequest struct {
    LotID       int64     `json:"lot_id" validate:"required,gt=0"`
    StartPrice  float64   `json:"start_price" validate:"required,gt=0"`
    MinStep     float64   `json:"min_step" validate:"required,gt=0"`
    StartTime   time.Time `json:"start_time" validate:"required"`
    EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
}

type UpdateAuctionRequest struct {
    StartPrice *float64
    MinStep    *float64
    StartTime  *time.Time
    EndTime    *time.Time
    Status     *string
}
