package handler

import (
    "context"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/timestamppb"
    
    pb "auction-system/pkg/api"
    "auction-system/internal/application/dto/user"
    userUseCase "auction-system/internal/application/usecase/user"
)

type UserHandler struct {
    pb.UnimplementedUserServiceServer
    CreateUserUC     userUseCase.CreateUserUseCaseInterface
    GetUserUC        userUseCase.GetUserUseCaseInterface
    UpdateUserUC     userUseCase.UpdateUserUseCaseInterface
    DeleteUserUC     userUseCase.DeleteUserUseCaseInterface
    GetAllUsersUC    userUseCase.GetAllUsersUseCaseInterface
    UpdateBalanceUC  userUseCase.UpdateBalanceUseCaseInterface
}

func NewUserHandler(
    createUserUC *userUseCase.CreateUserUseCase,
    getUserUC *userUseCase.GetUserUseCase,
    updateUserUC *userUseCase.UpdateUserUseCase,
    deleteUserUC *userUseCase.DeleteUserUseCase,
    getAllUsersUC *userUseCase.GetAllUserUseCase,
    updateBalanceUC *userUseCase.UpdateBalanceUseCase,
) *UserHandler {
    return &UserHandler{
        CreateUserUC:    createUserUC,
        GetUserUC:       getUserUC,
        UpdateUserUC:    updateUserUC,
        DeleteUserUC:    deleteUserUC,
        GetAllUsersUC:   getAllUsersUC,
        UpdateBalanceUC: updateBalanceUC,
    }
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    createReq := &user.CreateUserRequest{
        Username: req.Username,
        Email:    req.Email,
    }
    
    resp, err := h.CreateUserUC.Execute(ctx, createReq)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.CreateUserResponse{
        User: toProtoUser(resp),
    }, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    resp, err := h.GetUserUC.Execute(ctx, req.Id)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.GetUserResponse{
        User: toProtoUser(resp),
    }, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
    updateReq := &user.UpdateUserRequest{
        Username: req.Username,
        Email:    req.Email,
    }
    
    resp, err := h.UpdateUserUC.Execute(ctx, req.Id, updateReq)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.UpdateUserResponse{
        User: toProtoUser(resp),
    }, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
    err := h.DeleteUserUC.Execute(ctx, req.Id)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.DeleteUserResponse{}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
    resp, err := h.GetAllUsersUC.Execute(ctx)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    users := make([]*pb.User, 0, len(resp.Users))
    for _, u := range resp.Users {
        users = append(users, toProtoUser(&u))
    }

    return &pb.ListUsersResponse{
        Users:      users,
        TotalCount: int32(resp.Total),
    }, nil
}

func (h *UserHandler) UpdateBalance(ctx context.Context, req *pb.UpdateBalanceRequest) (*pb.UpdateBalanceResponse, error) {
    input := userUseCase.UpdateBalanceInput{
        UserID: req.UserId,
        Amount: req.Amount,
    }
    
    err := h.UpdateBalanceUC.Execute(ctx, input)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    updatedUser, err := h.GetUserUC.Execute(ctx, req.UserId)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &pb.UpdateBalanceResponse{
        User: toProtoUser(updatedUser),
    }, nil
}

func toProtoUser(u *user.UserResponse) *pb.User {
    return &pb.User{
        Id:        u.ID,
        Username:  u.Username,
        Email:     u.Email,
        Balance:   u.Balance,
        CreatedAt: timestamppb.New(u.CreatedAt),
        UpdatedAt: timestamppb.New(u.CreatedAt),
    }
}
