package entity

import (
    "time"
)

type Lot struct {
    ID          int64     `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    StartPrice  float64   `json:"start_price"`
    CreatorID   int64     `json:"creator_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
