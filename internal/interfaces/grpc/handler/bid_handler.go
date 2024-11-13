package handler

import (
    "context"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    pb "auction-system/pkg/api"
    bidUseCase "auction-system/internal/application/usecase/bid"
    "auction-system/internal/application/dto/bid"
    "google.golang.org/protobuf/types/known/timestamppb"
)

type BidHandler struct {
    pb.UnimplementedBidServiceServer
    placeBidUseCase *bidUseCase.PlaceBidUseCase
    getBidUseCase   *bidUseCase.GetBidUseCase
    listBidsUseCase *bidUseCase.ListBidsUseCase
}

func NewBidHandler(
    placeBidUseCase *bidUseCase.PlaceBidUseCase,
    getBidUseCase *bidUseCase.GetBidUseCase,
    listBidsUseCase *bidUseCase.ListBidsUseCase,
) *BidHandler {
    return &BidHandler{
        placeBidUseCase: placeBidUseCase,
        getBidUseCase:   getBidUseCase,
        listBidsUseCase: listBidsUseCase,
    }
}

func (h *BidHandler) PlaceBid(ctx context.Context, req *pb.PlaceBidRequest) (*pb.PlaceBidResponse, error) {
    placeBidReq := &bid.PlaceBidRequest{
        AuctionID: req.AuctionId,
        UserID:    req.UserId,
        Amount:    req.Amount,
    }

    result, err := h.placeBidUseCase.Execute(ctx, placeBidReq)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.PlaceBidResponse{
        Bid: mapBidToProto(result),
    }, nil
}

func (h *BidHandler) GetBid(ctx context.Context, req *pb.GetBidRequest) (*pb.GetBidResponse, error) {
    result, err := h.getBidUseCase.Execute(ctx, req.Id)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.GetBidResponse{
        Bid: mapBidToProto(result),
    }, nil
}

func (h *BidHandler) ListBids(ctx context.Context, req *pb.ListBidsRequest) (*pb.ListBidsResponse, error) {
    listBidsReq := &bid.ListBidsRequest{
        AuctionID:  req.AuctionId,
        PageSize:   int(req.PageSize),
        PageNumber: int(req.PageNumber),
    }

    result, err := h.listBidsUseCase.Execute(ctx, listBidsReq)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    bids := make([]*pb.Bid, len(result.Bids))
    for i, b := range result.Bids {
        bids[i] = mapBidToProto(&b)
    }

    return &pb.ListBidsResponse{
        Bids:       bids,
        TotalCount: result.TotalCount,
    }, nil
}

func mapBidToProto(b *bid.BidResponse) *pb.Bid {
    return &pb.Bid{
        Id:        b.ID,
        AuctionId: b.AuctionID,
        UserId:    b.UserID,
        Amount:    b.Amount,
        CreatedAt: timestamppb.New(b.CreatedAt),
        UpdatedAt: timestamppb.New(b.UpdatedAt),
    }
}
