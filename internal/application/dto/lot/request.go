package lot

type CreateLotRequest struct {
    Title       string  `json:"title" validate:"required,min=3,max=100"`
    Description string  `json:"description" validate:"required"`
    StartPrice  float64 `json:"start_price" validate:"required,gt=0"`
    CreatorID   int64   `json:"creator_id" validate:"required,gt=0"`
}

type UpdateLotRequest struct {
    Title       string  `json:"title,omitempty" validate:"omitempty,min=3,max=100"`
    Description string  `json:"description,omitempty"`
    StartPrice  float64 `json:"start_price,omitempty" validate:"omitempty,gt=0"`
}
