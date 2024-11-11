package app

import (
    "context"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    
    "auction-system/internal/config"
    "auction-system/internal/interfaces/grpc/handler"
    "auction-system/internal/application/usecase/user"
    pb "auction-system/pkg/api"
)

type App struct {
    cfg         *config.Config
    grpcServer  *grpc.Server
    httpServer  *http.Server
    userHandler *handler.UserHandler
}

func NewApp(cfg *config.Config, 
    createUserUC *user.CreateUserUseCase,
    getUserUC *user.GetUserUseCase,
    updateUserUC *user.UpdateUserUseCase,
    deleteUserUC *user.DeleteUserUseCase,
    getAllUsersUC *user.GetAllUserUseCase,
    updateBalanceUC *user.UpdateBalanceUseCase,
) *App {
    userHandler := handler.NewUserHandler(
        createUserUC,
        getUserUC,
        updateUserUC,
        deleteUserUC,
        getAllUsersUC,
        updateBalanceUC,
    )

    return &App{
        cfg:         cfg,
        userHandler: userHandler,
    }
}

func (a *App) Run(ctx context.Context) error {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    go func() {
        if err := a.runGRPCServer(); err != nil {
            log.Printf("Failed to run gRPC server: %v", err)
            cancel()
        }
    }()

    go func() {
        if err := a.runHTTPServer(ctx); err != nil && err != http.ErrServerClosed {
            log.Printf("Failed to run HTTP server: %v", err)
            cancel()
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    select {
    case <-ctx.Done():
        return ctx.Err()
    case sig := <-quit:
        log.Printf("Received signal: %v", sig)
    }

    return a.shutdown(ctx)
}

func (a *App) runGRPCServer() error {
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port))
    if err != nil {
        return fmt.Errorf("failed to listen: %v", err)
    }

    a.grpcServer = grpc.NewServer()
    pb.RegisterUserServiceServer(a.grpcServer, a.userHandler)

    log.Printf("Starting gRPC server on %s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port)
    return a.grpcServer.Serve(listener)
}

func (a *App) runHTTPServer(ctx context.Context) error {
    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

    endpoint := fmt.Sprintf("%s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port)
    if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
        return fmt.Errorf("failed to register gateway: %v", err)
    }

    a.httpServer = &http.Server{
        Addr:         fmt.Sprintf("%s:%d", a.cfg.HTTP.Host, a.cfg.HTTP.Port),
        Handler:      mux,
        ReadTimeout:  a.cfg.Server.ReadTimeout,
        WriteTimeout: a.cfg.Server.WriteTimeout,
    }

    log.Printf("Starting HTTP server on %s:%d", a.cfg.HTTP.Host, a.cfg.HTTP.Port)
    return a.httpServer.ListenAndServe()
}

func (a *App) shutdown(ctx context.Context) error {
    if a.httpServer != nil {
        if err := a.httpServer.Shutdown(ctx); err != nil {
            log.Printf("Failed to shutdown HTTP server: %v", err)
        }
    }

    if a.grpcServer != nil {
        a.grpcServer.GracefulStop()
    }

    return nil
}
