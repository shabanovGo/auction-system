package tests

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "google.golang.org/protobuf/types/known/timestamppb"
    
    pb "auction-system/pkg/api"
    dto "auction-system/internal/application/dto/auction"
    "auction-system/internal/domain/entity"
    "auction-system/internal/interfaces/grpc/handler"
    auctionUC "auction-system/internal/application/usecase/auction"
)

type mockCreateAuctionUC struct {
    mock.Mock
}

func (m *mockCreateAuctionUC) Execute(ctx context.Context, req *dto.CreateAuctionRequest) (*dto.AuctionResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.AuctionResponse), args.Error(1)
}

type mockGetAuctionUC struct {
    mock.Mock
}

func (m *mockGetAuctionUC) Execute(ctx context.Context, id int64) (*dto.AuctionResponse, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.AuctionResponse), args.Error(1)
}

type mockUpdateAuctionUC struct {
    mock.Mock
}

func (m *mockUpdateAuctionUC) Execute(ctx context.Context, id int64, req *dto.UpdateAuctionRequest) (*dto.AuctionResponse, error) {
    args := m.Called(ctx, id, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.AuctionResponse), args.Error(1)
}

type mockListAuctionsUC struct {
    mock.Mock
}

func (m *mockListAuctionsUC) Execute(ctx context.Context, page, pageSize int, status *string) (*dto.ListAuctionsResponse, error) {
    args := m.Called(ctx, page, pageSize, status)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.ListAuctionsResponse), args.Error(1)
}

func createTestAuctionResponse() *dto.AuctionResponse {
    winnerID := int64(1)
    winnerBidID := int64(1)
    now := time.Now()
    return &dto.AuctionResponse{
        ID:           1,
        LotID:        1,
        StartPrice:   100,
        MinStep:      10,
        CurrentPrice: 110,
        StartTime:    now,
        EndTime:      now.Add(24 * time.Hour),
        Status:       entity.AuctionStatusActive,
        WinnerID:     &winnerID,
        WinnerBidID:  &winnerBidID,
        CreatedAt:    now,
        UpdatedAt:    now,
    }
}

var _ auctionUC.CreateAuctionUseCaseInterface = (*mockCreateAuctionUC)(nil)
var _ auctionUC.GetAuctionUseCaseInterface = (*mockGetAuctionUC)(nil)
var _ auctionUC.UpdateAuctionUseCaseInterface = (*mockUpdateAuctionUC)(nil)
var _ auctionUC.ListAuctionsUseCaseInterface = (*mockListAuctionsUC)(nil)

func TestCreateAuction(t *testing.T) {
    mockUC := new(mockCreateAuctionUC)
    h := &handler.AuctionHandler{
        CreateAuctionUC: mockUC,
    }

    ctx := context.Background()
    now := time.Now()
    req := &pb.CreateAuctionRequest{
        LotId:      1,
        StartPrice: 100,
        MinStep:    10,
        StartTime:  timestamppb.New(now),
        EndTime:    timestamppb.New(now.Add(24 * time.Hour)),
    }

    expectedResp := createTestAuctionResponse()
    mockUC.On("Execute", ctx, mock.MatchedBy(func(req *dto.CreateAuctionRequest) bool {
        return req.LotID == 1 && req.StartPrice == 100 && req.MinStep == 10
    })).Return(expectedResp, nil)

    resp, err := h.CreateAuction(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.ID, resp.Auction.Id)
    mockUC.AssertExpectations(t)
}

func TestGetAuction(t *testing.T) {
    mockUC := new(mockGetAuctionUC)
    h := &handler.AuctionHandler{
        GetAuctionUC: mockUC,
    }

    ctx := context.Background()
    req := &pb.GetAuctionRequest{Id: 1}

    expectedResp := createTestAuctionResponse()
    mockUC.On("Execute", ctx, int64(1)).Return(expectedResp, nil)

    resp, err := h.GetAuction(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.ID, resp.Auction.Id)
    mockUC.AssertExpectations(t)
}

func TestUpdateAuction(t *testing.T) {
    mockUC := new(mockUpdateAuctionUC)
    h := &handler.AuctionHandler{
        UpdateAuctionUC: mockUC,
    }

    ctx := context.Background()
    now := time.Now()
    req := &pb.UpdateAuctionRequest{
        Id:         1,
        StartPrice: 100,
        MinStep:    10,
        StartTime:  timestamppb.New(now),
        EndTime:    timestamppb.New(now.Add(24 * time.Hour)),
        Status:     string(entity.AuctionStatusActive),
    }

    expectedResp := createTestAuctionResponse()
    mockUC.On("Execute", ctx, int64(1), mock.MatchedBy(func(req *dto.UpdateAuctionRequest) bool {
        return *req.StartPrice == 100 && *req.MinStep == 10 && *req.Status == string(entity.AuctionStatusActive)
    })).Return(expectedResp, nil)

    resp, err := h.UpdateAuction(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.ID, resp.Auction.Id)
    mockUC.AssertExpectations(t)
}

func TestListAuctions(t *testing.T) {
    mockUC := new(mockListAuctionsUC)
    h := &handler.AuctionHandler{
        ListAuctionsUC: mockUC,
    }

    ctx := context.Background()
    status := "ACTIVE"
    req := &pb.ListAuctionsRequest{
        PageSize:   10,
        PageNumber: 1,
        Status:     &status,
    }

    expectedResp := &dto.ListAuctionsResponse{
        Auctions:   []dto.AuctionResponse{*createTestAuctionResponse()},
        TotalCount: 1,
    }
    mockUC.On("Execute", ctx, 1, 10, &status).Return(expectedResp, nil)

    resp, err := h.ListAuctions(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, int32(1), resp.TotalCount)
    assert.Len(t, resp.Auctions, 1)
    mockUC.AssertExpectations(t)
}

func TestMapAuctionToProto(t *testing.T) {
    auctionResp := createTestAuctionResponse()
    protoAuction := handler.MapAuctionToProto(auctionResp)

    assert.NotNil(t, protoAuction)
    assert.Equal(t, auctionResp.ID, protoAuction.Id)
    assert.Equal(t, auctionResp.LotID, protoAuction.LotId)
    assert.Equal(t, auctionResp.StartPrice, protoAuction.StartPrice)
    assert.Equal(t, auctionResp.MinStep, protoAuction.MinStep)
    assert.Equal(t, auctionResp.CurrentPrice, protoAuction.CurrentPrice)
    assert.Equal(t, string(auctionResp.Status), protoAuction.Status)
    assert.Equal(t, *auctionResp.WinnerID, protoAuction.WinnerId)
    assert.Equal(t, *auctionResp.WinnerBidID, protoAuction.WinnerBidId)
}
