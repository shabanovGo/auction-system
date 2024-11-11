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
    createLotUC *lotUseCase.CreateLotUseCase
    getLotUC    *lotUseCase.GetLotUseCase
    updateLotUC *lotUseCase.UpdateLotUseCase
    deleteLotUC *lotUseCase.DeleteLotUseCase
    getLotsUC   *lotUseCase.GetLotsUseCase
}

func NewLotHandler(
    createLotUC *lotUseCase.CreateLotUseCase,
    getLotUC *lotUseCase.GetLotUseCase,
    updateLotUC *lotUseCase.UpdateLotUseCase,
    deleteLotUC *lotUseCase.DeleteLotUseCase,
    getLotsUC *lotUseCase.GetLotsUseCase,
) *LotHandler {
    return &LotHandler{
        createLotUC: createLotUC,
        getLotUC:    getLotUC,
        updateLotUC: updateLotUC,
        deleteLotUC: deleteLotUC,
        getLotsUC:   getLotsUC,
    }
}

func (h *LotHandler) CreateLot(ctx context.Context, req *pb.CreateLotRequest) (*pb.CreateLotResponse, error) {
    createReq := &lot.CreateLotRequest{
        Title:       req.Title,
        Description: req.Description,
        StartPrice:  req.StartPrice,
        CreatorID:   req.CreatorId,
    }

    lotResp, err := h.createLotUC.Execute(ctx, createReq)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.CreateLotResponse{
        Lot: toLotProto(lotResp),
    }, nil
}

func (h *LotHandler) GetLot(ctx context.Context, req *pb.GetLotRequest) (*pb.GetLotResponse, error) {
    lotResp, err := h.getLotUC.Execute(ctx, req.Id)
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

    lotResp, err := h.updateLotUC.Execute(ctx, req.Id, updateReq)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.UpdateLotResponse{
        Lot: toLotProto(lotResp),
    }, nil
}

func (h *LotHandler) DeleteLot(ctx context.Context, req *pb.DeleteLotRequest) (*pb.DeleteLotResponse, error) {
    err := h.deleteLotUC.Execute(ctx, req.Id)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.DeleteLotResponse{
        Success: true,
    }, nil
}

func (h *LotHandler) ListLots(ctx context.Context, req *pb.ListLotsRequest) (*pb.ListLotsResponse, error) {
    lotsResp, err := h.getLotsUC.Execute(ctx, int(req.PageNumber), int(req.PageSize))
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
