package auction

import (
    "context"
    dto "auction-system/internal/application/dto/auction"
)

type CreateAuctionUseCaseInterface interface {
    Execute(ctx context.Context, req *dto.CreateAuctionRequest) (*dto.AuctionResponse, error)
}

type GetAuctionUseCaseInterface interface {
    Execute(ctx context.Context, id int64) (*dto.AuctionResponse, error)
}

type UpdateAuctionUseCaseInterface interface {
    Execute(ctx context.Context, id int64, req *dto.UpdateAuctionRequest) (*dto.AuctionResponse, error)
}

type ListAuctionsUseCaseInterface interface {
    Execute(ctx context.Context, page, pageSize int, status *string) (*dto.ListAuctionsResponse, error)
}
