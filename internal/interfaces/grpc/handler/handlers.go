package handler

import (
    "context"
    "fmt"
    "net"
    "net/http"
    "auction-system/internal/config"
    "auction-system/pkg/api"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
)

type Handlers struct {
    userHandler    *UserHandler
    auctionHandler *AuctionHandler
    lotHandler     *LotHandler
    bidHandler     *BidHandler
    grpcServer     *grpc.Server
    httpServer     *http.Server
}

func NewHandlers(userHandler *UserHandler, auctionHandler *AuctionHandler, lotHandler *LotHandler, bidHandler *BidHandler) *Handlers {
    return &Handlers{
        userHandler:    userHandler,
        auctionHandler: auctionHandler,
        lotHandler:     lotHandler,
        bidHandler:     bidHandler,
    }
}

func (h *Handlers) Serve(ctx context.Context, cfg *config.Config) error {
    grpcServer := grpc.NewServer()
    h.grpcServer = grpcServer

    api.RegisterUserServiceServer(grpcServer, h.userHandler)
    api.RegisterAuctionServiceServer(grpcServer, h.auctionHandler)
    api.RegisterLotServiceServer(grpcServer, h.lotHandler)
    api.RegisterBidServiceServer(grpcServer, h.bidHandler)

    grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port))
    if err != nil {
        return fmt.Errorf("failed to listen grpc: %v", err)
    }

    go func() {
        if err := grpcServer.Serve(grpcListener); err != nil {
            fmt.Printf("failed to serve grpc: %v\n", err)
        }
    }()

    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithInsecure()}

    if err := api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port), opts); err != nil {
        return fmt.Errorf("failed to register user service handler: %v", err)
    }
    if err := api.RegisterAuctionServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port), opts); err != nil {
        return fmt.Errorf("failed to register auction service handler: %v", err)
    }
    if err := api.RegisterLotServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port), opts); err != nil {
        return fmt.Errorf("failed to register lot service handler: %v", err)
    }
    if err := api.RegisterBidServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port), opts); err != nil {
        return fmt.Errorf("failed to register bid service handler: %v", err)
    }

    h.httpServer = &http.Server{
        Addr:    fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
        Handler: mux,
    }

    if err := h.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        return fmt.Errorf("failed to serve http: %v", err)
    }

    return nil
}

func (h *Handlers) Shutdown(ctx context.Context) error {
    if err := h.httpServer.Shutdown(ctx); err != nil {
        return fmt.Errorf("failed to shutdown http server: %v", err)
    }

    h.grpcServer.GracefulStop()

    return nil
}
