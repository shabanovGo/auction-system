package app

import (
    "context"
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "auction-system/internal/config"
    "google.golang.org/grpc"
)

type App struct {
    cfg        *config.Config
    grpcServer *grpc.Server
    done       chan struct{}
    wg         sync.WaitGroup
}

func NewApp(cfg *config.Config) *App {
    return &App{
        cfg:  cfg,
        done: make(chan struct{}),
    }
}

func (a *App) Run(ctx context.Context) error {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    a.grpcServer = grpc.NewServer()

    grpcAddr := fmt.Sprintf("%s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port)
    grpcListener, err := net.Listen("tcp", grpcAddr)
    if err != nil {
        return fmt.Errorf("failed to listen grpc: %w", err)
    }

    a.wg.Add(1)
    go func() {
        defer a.wg.Done()
        log.Printf("Starting gRPC server on %s", grpcAddr)
        if err := a.grpcServer.Serve(grpcListener); err != nil {
            log.Printf("gRPC server error: %v", err)
            cancel()
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

    select {
    case <-ctx.Done():
        log.Println("Context cancelled")
    case sig := <-quit:
        log.Printf("Received signal: %v", sig)
    }

    return a.shutdown(ctx)
}

func (a *App) shutdown(ctx context.Context) error {
    shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    a.grpcServer.GracefulStop()

    close(a.done)

    waitCh := make(chan struct{})
    go func() {
        a.wg.Wait()
        close(waitCh)
    }()

    select {
    case <-waitCh:
        log.Println("Gracefully shut down")
        return nil
    case <-shutdownCtx.Done():
        return fmt.Errorf("shutdown timeout exceeded")
    }
}
