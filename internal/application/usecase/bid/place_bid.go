package bid

import (
    "context"
    "auction-system/internal/application/dto/bid"
    "auction-system/internal/domain/repository"
    "auction-system/internal/domain/errors"
    "auction-system/internal/domain/entity"
)

type PlaceBidUseCase struct {
    bidRepo     repository.BidRepository
    auctionRepo repository.AuctionRepository
    userRepo    repository.UserRepository
}

func NewPlaceBidUseCase(
    bidRepo repository.BidRepository,
    auctionRepo repository.AuctionRepository,
    userRepo repository.UserRepository,
) *PlaceBidUseCase {
    return &PlaceBidUseCase{
        bidRepo:     bidRepo,
        auctionRepo: auctionRepo,
        userRepo:    userRepo,
    }
}

func (uc *PlaceBidUseCase) Execute(ctx context.Context, req *bid.PlaceBidRequest) (*bid.BidResponse, error) {
    auction, err := uc.auctionRepo.GetByID(ctx, req.AuctionID)
    if err != nil {
        return nil, err
    }

    if auction.Status != entity.AuctionStatusActive {
        return nil, errors.New(errors.ErrorTypeValidation, "auction is not active", nil)
    }

    user, err := uc.userRepo.GetByID(ctx, req.UserID)
    if err != nil {
        return nil, err
    }

    if user.Balance < req.Amount {
        return nil, errors.New(errors.ErrorTypeValidation, "insufficient funds", nil)
    }

    if auction.CurrentPrice + auction.MinStep > req.Amount {
        return nil, errors.New(errors.ErrorTypeValidation, "bid amount must be greater than current price plus minimum step", nil)
    }

    bidEntity := bid.ToEntity(req)
    
    if err := uc.bidRepo.Create(ctx, bidEntity); err != nil {
        return nil, err
    }

    auction.CurrentPrice = req.Amount
    _, err = uc.auctionRepo.Update(ctx, auction.ID, auction)
    if err != nil {
        return nil, err
    }

    return bid.FromEntity(bidEntity), nil
}
