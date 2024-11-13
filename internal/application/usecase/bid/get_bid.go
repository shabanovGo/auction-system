package bid

import (
    "context"
    "auction-system/internal/application/dto/bid"
    "auction-system/internal/domain/repository"
)

type GetBidUseCase struct {
    bidRepo repository.BidRepository
}

func NewGetBidUseCase(bidRepo repository.BidRepository) *GetBidUseCase {
    return &GetBidUseCase{
        bidRepo: bidRepo,
    }
}

func (uc *GetBidUseCase) Execute(ctx context.Context, id int64) (*bid.BidResponse, error) {
    bidEntity, err := uc.bidRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    return bid.FromEntity(bidEntity), nil
}
