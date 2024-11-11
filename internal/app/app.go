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
    pb "auction-system/pkg/api"
)

type App struct {
    cfg      *config.Config
    handlers *handler.Handlers
    grpcServer *grpc.Server
    httpServer *http.Server
}

func NewApp(
    cfg *config.Config,
    handlers *handler.Handlers,
) *App {
    return &App{
        cfg:      cfg,
        handlers: handlers,
    }
}

func (a *App) Run(ctx context.Context) error {
    // Настройка gRPC сервера
    a.grpcServer = grpc.NewServer()
    pb.RegisterUserServiceServer(a.grpcServer, a.handlers.UserHandler)
    pb.RegisterLotServiceServer(a.grpcServer, a.handlers.LotHandler)

    // Запуск gRPC сервера
    grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port))
    if err != nil {
        return fmt.Errorf("failed to listen grpc: %v", err)
    }

    go func() {
        log.Printf("Starting gRPC server on %s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port)
        if err := a.grpcServer.Serve(grpcListener); err != nil {
            log.Fatalf("Failed to serve gRPC: %v", err)
        }
    }()

    // Настройка HTTP сервера с gRPC-Gateway
    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
    
    endpoint := fmt.Sprintf("%s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port)
    
    if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
        return fmt.Errorf("failed to register user service handler: %v", err)
    }
    
    if err := pb.RegisterLotServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
        return fmt.Errorf("failed to register lot service handler: %v", err)
    }

    a.httpServer = &http.Server{
        Addr:    fmt.Sprintf("%s:%d", a.cfg.HTTP.Host, a.cfg.HTTP.Port),
        Handler: mux,
    }

    // Запуск HTTP сервера
    go func() {
        log.Printf("Starting HTTP server on %s:%d", a.cfg.HTTP.Host, a.cfg.HTTP.Port)
        if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to serve HTTP: %v", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down servers...")

    a.grpcServer.GracefulStop()
    if err := a.httpServer.Shutdown(ctx); err != nil {
        log.Printf("HTTP server shutdown error: %v", err)
    }

    return nil
}
