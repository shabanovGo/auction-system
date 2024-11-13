package tests

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    pb "auction-system/pkg/api"
    dto "auction-system/internal/application/dto/lot"
    "auction-system/internal/interfaces/grpc/handler"
    lotUC "auction-system/internal/application/usecase/lot"
)

type mockCreateLotUC struct {
    mock.Mock
}

func (m *mockCreateLotUC) Execute(ctx context.Context, req *dto.CreateLotRequest) (*dto.LotResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.LotResponse), args.Error(1)
}

type mockGetLotUC struct {
    mock.Mock
}

func (m *mockGetLotUC) Execute(ctx context.Context, id int64) (*dto.LotResponse, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.LotResponse), args.Error(1)
}

type mockUpdateLotUC struct {
    mock.Mock
}

func (m *mockUpdateLotUC) Execute(ctx context.Context, id int64, req *dto.UpdateLotRequest) (*dto.LotResponse, error) {
    args := m.Called(ctx, id, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.LotResponse), args.Error(1)
}

type mockDeleteLotUC struct {
    mock.Mock
}

func (m *mockDeleteLotUC) Execute(ctx context.Context, id int64) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

type mockListLotsUC struct {
    mock.Mock
}

func (m *mockListLotsUC) Execute(ctx context.Context, page, pageSize int) (*dto.LotListResponse, error) {
    args := m.Called(ctx, page, pageSize)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.LotListResponse), args.Error(1)
}

func createTestLotResponse() *dto.LotResponse {
    now := time.Now()
    return &dto.LotResponse{
        ID:          1,
        Title:       "Test Lot",
        Description: "Test Description",
        StartPrice:  100.0,
        CreatorID:   1,
        CreatedAt:   now,
    }
}

var _ lotUC.CreateLotUseCaseInterface = (*mockCreateLotUC)(nil)
var _ lotUC.GetLotUseCaseInterface = (*mockGetLotUC)(nil)
var _ lotUC.UpdateLotUseCaseInterface = (*mockUpdateLotUC)(nil)
var _ lotUC.DeleteLotUseCaseInterface = (*mockDeleteLotUC)(nil)
var _ lotUC.ListLotsUseCaseInterface = (*mockListLotsUC)(nil)

func TestCreateLot(t *testing.T) {
    mockUC := new(mockCreateLotUC)
    h := &handler.LotHandler{
        CreateLotUC: mockUC,
    }

    ctx := context.Background()
    req := &pb.CreateLotRequest{
        Title:       "Test Lot",
        Description: "Test Description",
        StartPrice:  100.0,
        CreatorId:   1,
    }

    expectedResp := createTestLotResponse()
    mockUC.On("Execute", ctx, mock.MatchedBy(func(req *dto.CreateLotRequest) bool {
        return req.Title == "Test Lot" && 
               req.Description == "Test Description" && 
               req.StartPrice == 100.0 &&
               req.CreatorID == 1
    })).Return(expectedResp, nil)

    resp, err := h.CreateLot(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.ID, resp.Lot.Id)
    mockUC.AssertExpectations(t)
}

func TestGetLot(t *testing.T) {
    mockUC := new(mockGetLotUC)
    h := &handler.LotHandler{
        GetLotUC: mockUC,
    }

    ctx := context.Background()
    req := &pb.GetLotRequest{Id: 1}

    expectedResp := createTestLotResponse()
    mockUC.On("Execute", ctx, int64(1)).Return(expectedResp, nil)

    resp, err := h.GetLot(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.ID, resp.Lot.Id)
    mockUC.AssertExpectations(t)
}

func TestUpdateLot(t *testing.T) {
    mockUC := new(mockUpdateLotUC)
    h := &handler.LotHandler{
        UpdateLotUC: mockUC,
    }

    ctx := context.Background()
    req := &pb.UpdateLotRequest{
        Id:          1,
        Title:       "Updated Lot",
        Description: "Updated Description",
        StartPrice:  150.0,
    }

    expectedResp := createTestLotResponse()
    expectedResp.Title = "Updated Lot"
    expectedResp.Description = "Updated Description"
    expectedResp.StartPrice = 150.0

    mockUC.On("Execute", ctx, int64(1), mock.MatchedBy(func(req *dto.UpdateLotRequest) bool {
        return req.Title == "Updated Lot" && 
               req.Description == "Updated Description" && 
               req.StartPrice == 150.0
    })).Return(expectedResp, nil)

    resp, err := h.UpdateLot(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.ID, resp.Lot.Id)
    assert.Equal(t, expectedResp.Title, resp.Lot.Title)
    assert.Equal(t, expectedResp.Description, resp.Lot.Description)
    assert.Equal(t, expectedResp.StartPrice, resp.Lot.StartPrice)
    mockUC.AssertExpectations(t)
}

func TestDeleteLot(t *testing.T) {
    mockUC := new(mockDeleteLotUC)
    h := &handler.LotHandler{
        DeleteLotUC: mockUC,
    }

    ctx := context.Background()
    req := &pb.DeleteLotRequest{Id: 1}

    mockUC.On("Execute", ctx, int64(1)).Return(nil)

    resp, err := h.DeleteLot(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.True(t, resp.Success)
    mockUC.AssertExpectations(t)
}

func TestListLots(t *testing.T) {
    mockUC := new(mockListLotsUC)
    h := &handler.LotHandler{
        ListLotsUC: mockUC,
    }

    ctx := context.Background()
    req := &pb.ListLotsRequest{
        PageSize:   10,
        PageNumber: 1,
    }

    expectedLots := []dto.LotResponse{*createTestLotResponse()}
    expectedResp := &dto.LotListResponse{
        Lots:  expectedLots,
        Total: 1,
    }
    
    mockUC.On("Execute", ctx, 1, 10).Return(expectedResp, nil)

    resp, err := h.ListLots(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, int32(expectedResp.Total), resp.TotalCount)
    assert.Len(t, resp.Lots, 1)
    assert.Equal(t, expectedLots[0].ID, resp.Lots[0].Id)
    assert.Equal(t, expectedLots[0].Title, resp.Lots[0].Title)
    assert.Equal(t, expectedLots[0].Description, resp.Lots[0].Description)
    mockUC.AssertExpectations(t)
}
