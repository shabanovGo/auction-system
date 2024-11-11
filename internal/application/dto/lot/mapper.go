package lot

import (
    "auction-system/internal/domain/entity"
)

func (r *CreateLotRequest) ToEntity() *entity.Lot {
    return &entity.Lot{
        Title:       r.Title,
        Description: r.Description,
        StartPrice:  r.StartPrice,
        CreatorID:   r.CreatorID,
    }
}

func FromEntity(lot *entity.Lot) *LotResponse {
    return &LotResponse{
        ID:          lot.ID,
        Title:       lot.Title,
        Description: lot.Description,
        StartPrice:  lot.StartPrice,
        CreatorID:   lot.CreatorID,
        CreatedAt:   lot.CreatedAt,
        UpdatedAt:   lot.UpdatedAt,
    }
}

func ToLotListResponse(lots []*entity.Lot, total int64) *LotListResponse {
    response := &LotListResponse{
        Lots:  make([]LotResponse, len(lots)),
        Total: total,
    }

    for i, lot := range lots {
        response.Lots[i] = *FromEntity(lot)
    }

    return response
}
