package bid

import (
    "context"
    "auction-system/internal/application/dto/bid"
    "auction-system/internal/domain/repository"
    "auction-system/internal/domain/errors"
)

type ListBidsUseCase struct {
    bidRepo repository.BidRepository
}

func NewListBidsUseCase(bidRepo repository.BidRepository) *ListBidsUseCase {
    return &ListBidsUseCase{
        bidRepo: bidRepo,
    }
}

func (uc *ListBidsUseCase) Execute(ctx context.Context, req *bid.ListBidsRequest) (*bid.ListBidsResponse, error) {
    if req.PageNumber < 1 {
        return nil, errors.New(errors.ErrorTypeValidation, "page number must be greater than 0", nil)
    }
    if req.PageSize < 1 {
        return nil, errors.New(errors.ErrorTypeValidation, "page size must be greater than 0", nil)
    }

    offset := (req.PageNumber - 1) * req.PageSize
    bids, total, err := uc.bidRepo.List(ctx, req.AuctionID, offset, req.PageSize)
    if err != nil {
        return nil, err
    }

    return &bid.ListBidsResponse{
        Bids:       bid.FromEntityList(bids),
        TotalCount: total,
    }, nil
}
