package tests

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    pb "auction-system/pkg/api"
    dto "auction-system/internal/application/dto/user"
    "auction-system/internal/interfaces/grpc/handler"
    userUC "auction-system/internal/application/usecase/user"
)

// Моки для юзкейсов
type mockCreateUserUC struct {
    mock.Mock
}

func (m *mockCreateUserUC) Execute(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.UserResponse), args.Error(1)
}

type mockGetUserUC struct {
    mock.Mock
}

func (m *mockGetUserUC) Execute(ctx context.Context, id int64) (*dto.UserResponse, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.UserResponse), args.Error(1)
}

type mockUpdateUserUC struct {
    mock.Mock
}

func (m *mockUpdateUserUC) Execute(ctx context.Context, id int64, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
    args := m.Called(ctx, id, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.UserResponse), args.Error(1)
}

type mockDeleteUserUC struct {
    mock.Mock
}

func (m *mockDeleteUserUC) Execute(ctx context.Context, id int64) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

type mockGetAllUsersUC struct {
    mock.Mock
}

func (m *mockGetAllUsersUC) Execute(ctx context.Context) (*dto.UserListResponse, error) {
    args := m.Called(ctx)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.UserListResponse), args.Error(1)
}

type mockUpdateBalanceUC struct {
    mock.Mock
}

func (m *mockUpdateBalanceUC) Execute(ctx context.Context, input userUC.UpdateBalanceInput) error {
    args := m.Called(ctx, input)
    return args.Error(0)
}

// Вспомогательные функции
func createTestUserResponse() *dto.UserResponse {
    now := time.Now()
    return &dto.UserResponse{
        ID:        1,
        Username:  "testuser",
        Email:     "test@example.com",
        Balance:   100.0,
        CreatedAt: now,
    }
}

// Проверка реализации интерфейсов
var _ userUC.CreateUserUseCaseInterface = (*mockCreateUserUC)(nil)
var _ userUC.GetUserUseCaseInterface = (*mockGetUserUC)(nil)
var _ userUC.UpdateUserUseCaseInterface = (*mockUpdateUserUC)(nil)
var _ userUC.DeleteUserUseCaseInterface = (*mockDeleteUserUC)(nil)
var _ userUC.GetAllUsersUseCaseInterface = (*mockGetAllUsersUC)(nil)
var _ userUC.UpdateBalanceUseCaseInterface = (*mockUpdateBalanceUC)(nil)

// Тесты
func TestCreateUser(t *testing.T) {
    mockUC := new(mockCreateUserUC)
    h := &handler.UserHandler{
        CreateUserUC: mockUC,
    }

    ctx := context.Background()
    req := &pb.CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
    }

    expectedResp := createTestUserResponse()
    mockUC.On("Execute", ctx, mock.MatchedBy(func(req *dto.CreateUserRequest) bool {
        return req.Username == "testuser" && req.Email == "test@example.com"
    })).Return(expectedResp, nil)

    resp, err := h.CreateUser(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.ID, resp.User.Id)
    assert.Equal(t, expectedResp.Username, resp.User.Username)
    assert.Equal(t, expectedResp.Email, resp.User.Email)
    mockUC.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
    mockUC := new(mockGetUserUC)
    h := &handler.UserHandler{
        GetUserUC: mockUC,
    }

    ctx := context.Background()
    req := &pb.GetUserRequest{Id: 1}

    expectedResp := createTestUserResponse()
    mockUC.On("Execute", ctx, int64(1)).Return(expectedResp, nil)

    resp, err := h.GetUser(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.ID, resp.User.Id)
    assert.Equal(t, expectedResp.Username, resp.User.Username)
    assert.Equal(t, expectedResp.Email, resp.User.Email)
    mockUC.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
    mockUC := new(mockUpdateUserUC)
    h := &handler.UserHandler{
        UpdateUserUC: mockUC,
    }

    ctx := context.Background()
    req := &pb.UpdateUserRequest{
        Id:       1,
        Username: "newusername",
        Email:    "newemail@example.com",
    }

    expectedResp := createTestUserResponse()
    expectedResp.Username = "newusername"
    expectedResp.Email = "newemail@example.com"

    mockUC.On("Execute", ctx, int64(1), mock.MatchedBy(func(req *dto.UpdateUserRequest) bool {
        return req.Username == "newusername" && req.Email == "newemail@example.com"
    })).Return(expectedResp, nil)

    resp, err := h.UpdateUser(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResp.ID, resp.User.Id)
    assert.Equal(t, expectedResp.Username, resp.User.Username)
    assert.Equal(t, expectedResp.Email, resp.User.Email)
    mockUC.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
    mockUC := new(mockDeleteUserUC)
    h := &handler.UserHandler{
        DeleteUserUC: mockUC,
    }

    ctx := context.Background()
    req := &pb.DeleteUserRequest{Id: 1}

    mockUC.On("Execute", ctx, int64(1)).Return(nil)

    resp, err := h.DeleteUser(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    mockUC.AssertExpectations(t)
}

func TestListUsers(t *testing.T) {
    mockUC := new(mockGetAllUsersUC)
    h := &handler.UserHandler{
        GetAllUsersUC: mockUC,
    }

    ctx := context.Background()
    req := &pb.ListUsersRequest{}

    expectedUsers := []dto.UserResponse{*createTestUserResponse()}
    expectedResp := &dto.UserListResponse{
        Users: expectedUsers,
        Total: 1,
    }
    
    mockUC.On("Execute", ctx).Return(expectedResp, nil)

    resp, err := h.ListUsers(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, int32(expectedResp.Total), resp.TotalCount)
    assert.Len(t, resp.Users, 1)
    assert.Equal(t, expectedUsers[0].ID, resp.Users[0].Id)
    assert.Equal(t, expectedUsers[0].Username, resp.Users[0].Username)
    assert.Equal(t, expectedUsers[0].Email, resp.Users[0].Email)
    mockUC.AssertExpectations(t)
}

func TestUpdateBalance(t *testing.T) {
    mockBalanceUC := new(mockUpdateBalanceUC)
    mockGetUserUC := new(mockGetUserUC)
    h := &handler.UserHandler{
        UpdateBalanceUC: mockBalanceUC,
        GetUserUC:      mockGetUserUC,
    }

    ctx := context.Background()
    req := &pb.UpdateBalanceRequest{
        UserId: 1,
        Amount: 50.0,
    }

    expectedUser := createTestUserResponse()
    expectedUser.Balance = 150.0

    mockBalanceUC.On("Execute", ctx, userUC.UpdateBalanceInput{
        UserID: req.UserId,
        Amount: req.Amount,
    }).Return(nil)

    mockGetUserUC.On("Execute", ctx, int64(1)).Return(expectedUser, nil)

    resp, err := h.UpdateBalance(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedUser.ID, resp.User.Id)
    assert.Equal(t, expectedUser.Username, resp.User.Username)
    assert.Equal(t, expectedUser.Email, resp.User.Email)
    assert.Equal(t, expectedUser.Balance, resp.User.Balance)
    mockBalanceUC.AssertExpectations(t)
    mockGetUserUC.AssertExpectations(t)
}
