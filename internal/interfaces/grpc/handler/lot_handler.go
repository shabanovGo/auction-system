package handler

import (
    "context"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/timestamppb"
    
    pb "auction-system/pkg/api"
    "auction-system/internal/application/dto/lot"
    lotUseCase "auction-system/internal/application/usecase/lot"
)

type LotHandler struct {
    pb.UnimplementedLotServiceServer
    CreateLotUC  lotUseCase.CreateLotUseCaseInterface
    GetLotUC     lotUseCase.GetLotUseCaseInterface
    UpdateLotUC  lotUseCase.UpdateLotUseCaseInterface
    DeleteLotUC  lotUseCase.DeleteLotUseCaseInterface
    ListLotsUC   lotUseCase.ListLotsUseCaseInterface
}

func NewLotHandler(
    createLotUC lotUseCase.CreateLotUseCaseInterface,
    getLotUC lotUseCase.GetLotUseCaseInterface,
    updateLotUC lotUseCase.UpdateLotUseCaseInterface,
    deleteLotUC lotUseCase.DeleteLotUseCaseInterface,
    listLotsUC lotUseCase.ListLotsUseCaseInterface,
) *LotHandler {
    return &LotHandler{
        CreateLotUC:  createLotUC,
        GetLotUC:     getLotUC,
        UpdateLotUC:  updateLotUC,
        DeleteLotUC:  deleteLotUC,
        ListLotsUC:   listLotsUC,
    }
}

func (h *LotHandler) CreateLot(ctx context.Context, req *pb.CreateLotRequest) (*pb.CreateLotResponse, error) {
    createReq := &lot.CreateLotRequest{
        Title:       req.Title,
        Description: req.Description,
        StartPrice:  req.StartPrice,
        CreatorID:   req.CreatorId,
    }

    lotResp, err := h.CreateLotUC.Execute(ctx, createReq)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.CreateLotResponse{
        Lot: toLotProto(lotResp),
    }, nil
}

func (h *LotHandler) GetLot(ctx context.Context, req *pb.GetLotRequest) (*pb.GetLotResponse, error) {
    lotResp, err := h.GetLotUC.Execute(ctx, req.Id)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.GetLotResponse{
        Lot: toLotProto(lotResp),
    }, nil
}

func (h *LotHandler) UpdateLot(ctx context.Context, req *pb.UpdateLotRequest) (*pb.UpdateLotResponse, error) {
    updateReq := &lot.UpdateLotRequest{
        Title:       req.Title,
        Description: req.Description,
        StartPrice:  req.StartPrice,
    }

    lotResp, err := h.UpdateLotUC.Execute(ctx, req.Id, updateReq)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.UpdateLotResponse{
        Lot: toLotProto(lotResp),
    }, nil
}

func (h *LotHandler) DeleteLot(ctx context.Context, req *pb.DeleteLotRequest) (*pb.DeleteLotResponse, error) {
    err := h.DeleteLotUC.Execute(ctx, req.Id)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.DeleteLotResponse{
        Success: true,
    }, nil
}

func (h *LotHandler) ListLots(ctx context.Context, req *pb.ListLotsRequest) (*pb.ListLotsResponse, error) {
    lotsResp, err := h.ListLotsUC.Execute(ctx, int(req.PageNumber), int(req.PageSize))
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    pbLots := make([]*pb.Lot, len(lotsResp.Lots))
    for i, lot := range lotsResp.Lots {
        pbLots[i] = toLotProto(&lot)
    }

    return &pb.ListLotsResponse{
        Lots:       pbLots,
        TotalCount: int32(lotsResp.Total),
    }, nil
}

func toLotProto(lot *lot.LotResponse) *pb.Lot {
    return &pb.Lot{
        Id:          lot.ID,
        Title:       lot.Title,
        Description: lot.Description,
        StartPrice:  lot.StartPrice,
        CreatorId:   lot.CreatorID,
        CreatedAt:   timestamppb.New(lot.CreatedAt),
        UpdatedAt:   timestamppb.New(lot.UpdatedAt),
    }
}
