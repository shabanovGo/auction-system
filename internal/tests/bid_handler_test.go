package tests

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    pb "auction-system/pkg/api"
    dto "auction-system/internal/application/dto/bid"
    "auction-system/internal/interfaces/grpc/handler"
    bidUC "auction-system/internal/application/usecase/bid"
)

type mockPlaceBidUC struct {
    mock.Mock
}

func (m *mockPlaceBidUC) Execute(ctx context.Context, req *dto.PlaceBidRequest) (*dto.BidResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.BidResponse), args.Error(1)
}

type mockGetBidUC struct {
    mock.Mock
}

func (m *mockGetBidUC) Execute(ctx context.Context, id int64) (*dto.BidResponse, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.BidResponse), args.Error(1)
}

type mockListBidsUC struct {
    mock.Mock
}

func (m *mockListBidsUC) Execute(ctx context.Context, req *dto.ListBidsRequest) (*dto.ListBidsResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.ListBidsResponse), args.Error(1)
}

func createTestBidResponse() *dto.BidResponse {
    now := time.Now()
    return &dto.BidResponse{
        ID:        1,
        AuctionID: 1,
        UserID:    1,
        Amount:    150.0,
        CreatedAt: now,
        UpdatedAt: now,
    }
}

var _ bidUC.PlaceBidUseCaseInterface = (*mockPlaceBidUC)(nil)
var _ bidUC.GetBidUseCaseInterface = (*mockGetBidUC)(nil)
var _ bidUC.ListBidsUseCaseInterface = (*mockListBidsUC)(nil)

func TestPlaceBid(t *testing.T) {
    mockUC := new(mockPlaceBidUC)
    h := &handler.BidHandler{
        PlaceBidUseCase: mockUC,
    }

    ctx := context.Background()
    req := &pb.PlaceBidRequest{
        AuctionId: 1,
        UserId:    1,
        Amount:    150.0,
    }

    expectedResp := createTestBidResponse()
    mockUC.On("Execute", ctx, mock.MatchedBy(func(req *dto.PlaceBidRequest) bool {
        return req.AuctionID == 1 && req.UserID == 1 && req.Amount == 150.0
    })).Return(expectedResp, nil)

    resp, err := h.PlaceBid(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.ID, resp.Bid.Id)
    assert.Equal(t, expectedResp.AuctionID, resp.Bid.AuctionId)
    assert.Equal(t, expectedResp.UserID, resp.Bid.UserId)
    assert.Equal(t, expectedResp.Amount, resp.Bid.Amount)
    mockUC.AssertExpectations(t)
}

func TestGetBid(t *testing.T) {
    mockUC := new(mockGetBidUC)
    h := &handler.BidHandler{
        GetBidUseCase: mockUC,
    }

    ctx := context.Background()
    req := &pb.GetBidRequest{Id: 1}

    expectedResp := createTestBidResponse()
    mockUC.On("Execute", ctx, int64(1)).Return(expectedResp, nil)

    resp, err := h.GetBid(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.ID, resp.Bid.Id)
    assert.Equal(t, expectedResp.AuctionID, resp.Bid.AuctionId)
    assert.Equal(t, expectedResp.UserID, resp.Bid.UserId)
    assert.Equal(t, expectedResp.Amount, resp.Bid.Amount)
    mockUC.AssertExpectations(t)
}

func TestListBids(t *testing.T) {
    mockUC := new(mockListBidsUC)
    h := &handler.BidHandler{
        ListBidsUseCase: mockUC,
    }

    ctx := context.Background()
    req := &pb.ListBidsRequest{
        AuctionId:  1,
        PageSize:   10,
        PageNumber: 1,
    }

    expectedBids := []dto.BidResponse{*createTestBidResponse()}
    expectedResp := &dto.ListBidsResponse{
        Bids:       expectedBids,
        TotalCount: int64(1),
    }
    
    mockUC.On("Execute", ctx, mock.MatchedBy(func(req *dto.ListBidsRequest) bool {
        return req.AuctionID == 1 && req.PageSize == 10 && req.PageNumber == 1
    })).Return(expectedResp, nil)

    resp, err := h.ListBids(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.TotalCount, int64(resp.TotalCount))
    assert.Len(t, resp.Bids, 1)
    assert.Equal(t, expectedBids[0].ID, resp.Bids[0].Id)
    assert.Equal(t, expectedBids[0].AuctionID, resp.Bids[0].AuctionId)
    assert.Equal(t, expectedBids[0].UserID, resp.Bids[0].UserId)
    assert.Equal(t, expectedBids[0].Amount, resp.Bids[0].Amount)
    mockUC.AssertExpectations(t)
}
