package bid

import (
    "context"
    dto "auction-system/internal/application/dto/bid"
)

type PlaceBidUseCaseInterface interface {
    Execute(ctx context.Context, req *dto.PlaceBidRequest) (*dto.BidResponse, error)
}

type GetBidUseCaseInterface interface {
    Execute(ctx context.Context, id int64) (*dto.BidResponse, error)
}

type ListBidsUseCaseInterface interface {
    Execute(ctx context.Context, req *dto.ListBidsRequest) (*dto.ListBidsResponse, error)
}
