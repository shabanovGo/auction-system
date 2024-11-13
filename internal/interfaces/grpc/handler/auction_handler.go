package handler

import (
    "context"
    "google.golang.org/grpc/codes"
    grpcStatus "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/timestamppb"
    
    pb "auction-system/pkg/api"
    dto "auction-system/internal/application/dto/auction"
    auctionUseCase "auction-system/internal/application/usecase/auction"
)

type AuctionHandler struct {
    pb.UnimplementedAuctionServiceServer
    CreateAuctionUC auctionUseCase.CreateAuctionUseCaseInterface
    GetAuctionUC    auctionUseCase.GetAuctionUseCaseInterface
    UpdateAuctionUC auctionUseCase.UpdateAuctionUseCaseInterface
    ListAuctionsUC  auctionUseCase.ListAuctionsUseCaseInterface
}

func NewAuctionHandler(
    createAuctionUC auctionUseCase.CreateAuctionUseCaseInterface,
    getAuctionUC auctionUseCase.GetAuctionUseCaseInterface,
    updateAuctionUC auctionUseCase.UpdateAuctionUseCaseInterface,
    listAuctionsUC auctionUseCase.ListAuctionsUseCaseInterface,
) *AuctionHandler {
    return &AuctionHandler{
        CreateAuctionUC: createAuctionUC,
        GetAuctionUC:    getAuctionUC,
        UpdateAuctionUC: updateAuctionUC,
        ListAuctionsUC:  listAuctionsUC,
    }
}

func (h *AuctionHandler) CreateAuction(ctx context.Context, req *pb.CreateAuctionRequest) (*pb.CreateAuctionResponse, error) {
    createReq := &dto.CreateAuctionRequest{
        LotID:      req.LotId,
        StartPrice: req.StartPrice,
        MinStep:    req.MinStep,
        StartTime:  req.StartTime.AsTime(),
        EndTime:    req.EndTime.AsTime(),
    }

    result, err := h.CreateAuctionUC.Execute(ctx, createReq)
    if err != nil {
        return nil, grpcStatus.Error(codes.Internal, err.Error())
    }

    return &pb.CreateAuctionResponse{
        Auction: MapAuctionToProto(result),
    }, nil
}

func (h *AuctionHandler) GetAuction(ctx context.Context, req *pb.GetAuctionRequest) (*pb.GetAuctionResponse, error) {
    result, err := h.GetAuctionUC.Execute(ctx, req.Id)
    if err != nil {
        return nil, grpcStatus.Error(codes.Internal, err.Error())
    }

    return &pb.GetAuctionResponse{
        Auction: MapAuctionToProto(result),
    }, nil
}

func (h *AuctionHandler) UpdateAuction(ctx context.Context, req *pb.UpdateAuctionRequest) (*pb.UpdateAuctionResponse, error) {
    startPrice := req.StartPrice
    minStep := req.MinStep
    auctionStatus := req.Status
    
    startTime := req.StartTime.AsTime()
    endTime := req.EndTime.AsTime()

    updateReq := &dto.UpdateAuctionRequest{
        StartPrice: &startPrice,
        MinStep:    &minStep,
        StartTime:  &startTime,
        EndTime:    &endTime,
        Status:     &auctionStatus,
    }

    result, err := h.UpdateAuctionUC.Execute(ctx, req.Id, updateReq)
    if err != nil {
        return nil, grpcStatus.Error(codes.Internal, err.Error())
    }

    return &pb.UpdateAuctionResponse{
        Auction: MapAuctionToProto(result),
    }, nil
}

func (h *AuctionHandler) ListAuctions(ctx context.Context, req *pb.ListAuctionsRequest) (*pb.ListAuctionsResponse, error) {
    var auctionStatus *string
    if req.Status != nil {
        s := *req.Status
        auctionStatus = &s
    }

    result, err := h.ListAuctionsUC.Execute(ctx, int(req.PageNumber), int(req.PageSize), auctionStatus)
    if err != nil {
        return nil, grpcStatus.Error(codes.Internal, err.Error())
    }

    auctions := make([]*pb.Auction, len(result.Auctions))
    for i, a := range result.Auctions {
        auctions[i] = MapAuctionToProto(&a)
    }

    return &pb.ListAuctionsResponse{
        Auctions:   auctions,
        TotalCount: int32(result.TotalCount),
    }, nil
}

func MapAuctionToProto(a *dto.AuctionResponse) *pb.Auction {
    var winnerID, winnerBidID int64
    if a.WinnerID != nil {
        winnerID = *a.WinnerID
    }
    if a.WinnerBidID != nil {
        winnerBidID = *a.WinnerBidID
    }

    return &pb.Auction{
        Id:           a.ID,
        LotId:        a.LotID,
        StartPrice:   a.StartPrice,
        MinStep:      a.MinStep,
        CurrentPrice: a.CurrentPrice,
        StartTime:    timestamppb.New(a.StartTime),
        EndTime:      timestamppb.New(a.EndTime),
        Status:       string(a.Status),
        WinnerId:     winnerID,
        WinnerBidId:  winnerBidID,
        CreatedAt:    timestamppb.New(a.CreatedAt),
        UpdatedAt:    timestamppb.New(a.UpdatedAt),
    }
}
