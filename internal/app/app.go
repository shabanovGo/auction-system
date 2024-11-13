package app

import (
    "context"
    "database/sql"
    _ "github.com/lib/pq"

    "auction-system/internal/config"
    "auction-system/internal/infrastructure/persistence/postgres"
    userUseCase "auction-system/internal/application/usecase/user"
    lotUseCase "auction-system/internal/application/usecase/lot"
    auctionUseCase "auction-system/internal/application/usecase/auction"
    bidUseCase "auction-system/internal/application/usecase/bid"
    handler "auction-system/internal/interfaces/grpc/handler"
    "auction-system/internal/worker"
    notificationDomain "auction-system/internal/domain/notification"
    notificationInfra "auction-system/internal/infrastructure/notification"
)

type App struct {
    cfg      *config.Config
    db       *sql.DB
    handlers *handler.Handlers
    worker   *worker.Worker
}

func NewApp(cfg *config.Config) (*App, error) {
    db, err := initDatabase(cfg)
    if err != nil {
        return nil, err
    }

    repos := initRepositories(db)
    services := initServices()
    useCases := initUseCases(repos)
    handlers := initHandlers(useCases)
    worker := initWorker(repos, services)

    return &App{
        cfg:      cfg,
        db:       db,
        handlers: handlers,
        worker:   worker,
    }, nil
}

func (a *App) Run(ctx context.Context) error {
    go a.worker.Start(ctx)

    return a.handlers.Serve(ctx, a.cfg)
}

func (a *App) Shutdown(ctx context.Context) error {
    if err := a.handlers.Shutdown(ctx); err != nil {
        return err
    }
    return a.db.Close()
}

func initDatabase(cfg *config.Config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.Database.GetDSN())
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}

type repositories struct {
    userRepo    *postgres.UserRepository
    lotRepo     *postgres.LotRepository
    auctionRepo *postgres.AuctionRepository
    bidRepo     *postgres.BidRepository
}

func initRepositories(db *sql.DB) *repositories {
    return &repositories{
        userRepo:    postgres.NewUserRepository(db),
        lotRepo:     postgres.NewLotRepository(db),
        auctionRepo: postgres.NewAuctionRepository(db),
        bidRepo:     postgres.NewBidRepository(db),
    }
}

type useCases struct {
    user    *userUseCases
    lot     *lotUseCases
    auction *auctionUseCases
    bid     *bidUseCases
}

type userUseCases struct {
    create        *userUseCase.CreateUserUseCase
    get          *userUseCase.GetUserUseCase
    update       *userUseCase.UpdateUserUseCase
    delete       *userUseCase.DeleteUserUseCase
    getAll       *userUseCase.GetAllUserUseCase
    updateBalance *userUseCase.UpdateBalanceUseCase
}

type lotUseCases struct {
    create  *lotUseCase.CreateLotUseCase
    get     *lotUseCase.GetLotUseCase
    update  *lotUseCase.UpdateLotUseCase
    delete  *lotUseCase.DeleteLotUseCase
    getAll  *lotUseCase.GetLotsUseCase
}

type auctionUseCases struct {
    create  *auctionUseCase.CreateAuctionUseCase
    get     *auctionUseCase.GetAuctionUseCase
    update  *auctionUseCase.UpdateAuctionUseCase
    list    *auctionUseCase.ListAuctionsUseCase
}

type bidUseCases struct {
    place   *bidUseCase.PlaceBidUseCase
    get     *bidUseCase.GetBidUseCase
    list    *bidUseCase.ListBidsUseCase
}

type services struct {
    notifier notificationDomain.NotificationService
}

func initServices() *services {
    return &services{
        notifier: notificationInfra.NewMockNotificationAdapter(),
    }
}

func initUseCases(repos *repositories) *useCases {
    return &useCases{
        user: &userUseCases{
            create:        userUseCase.NewCreateUserUseCase(repos.userRepo),
            get:          userUseCase.NewGetUserUseCase(repos.userRepo),
            update:       userUseCase.NewUpdateUserUseCase(repos.userRepo),
            delete:       userUseCase.NewDeleteUserUseCase(repos.userRepo),
            getAll:       userUseCase.NewGetAllUserUseCase(repos.userRepo),
            updateBalance: userUseCase.NewUpdateBalanceUseCase(repos.userRepo),
        },
        lot: &lotUseCases{
            create:  lotUseCase.NewCreateLotUseCase(repos.lotRepo),
            get:     lotUseCase.NewGetLotUseCase(repos.lotRepo),
            update:  lotUseCase.NewUpdateLotUseCase(repos.lotRepo),
            delete:  lotUseCase.NewDeleteLotUseCase(repos.lotRepo),
            getAll:  lotUseCase.NewGetLotsUseCase(repos.lotRepo),
        },
        auction: &auctionUseCases{
            create:  auctionUseCase.NewCreateAuctionUseCase(repos.auctionRepo, repos.lotRepo),
            get:     auctionUseCase.NewGetAuctionUseCase(repos.auctionRepo),
            update:  auctionUseCase.NewUpdateAuctionUseCase(repos.auctionRepo),
            list:    auctionUseCase.NewListAuctionsUseCase(repos.auctionRepo),
        },
        bid: &bidUseCases{
            place:   bidUseCase.NewPlaceBidUseCase(repos.bidRepo, repos.auctionRepo, repos.userRepo),
            get:     bidUseCase.NewGetBidUseCase(repos.bidRepo),
            list:    bidUseCase.NewListBidsUseCase(repos.bidRepo),
        },
    }
}

func initHandlers(uc *useCases) *handler.Handlers {
    userHandler := handler.NewUserHandler(
        uc.user.create,
        uc.user.get,
        uc.user.update,
        uc.user.delete,
        uc.user.getAll,
        uc.user.updateBalance,
    )

    lotHandler := handler.NewLotHandler(
        uc.lot.create,
        uc.lot.get,
        uc.lot.update,
        uc.lot.delete,
        uc.lot.getAll,
    )

    auctionHandler := handler.NewAuctionHandler(
        uc.auction.create,
        uc.auction.get,
        uc.auction.update,
        uc.auction.list,
    )

    bidHandler := handler.NewBidHandler(
        uc.bid.place,
        uc.bid.get,
        uc.bid.list,
    )

    return handler.NewHandlers(userHandler, auctionHandler, lotHandler, bidHandler)
}

func initWorker(repos *repositories, services *services) *worker.Worker {
    return worker.NewWorker(
        repos.auctionRepo,
        repos.lotRepo,
        repos.userRepo,
        repos.bidRepo,
        services.notifier,
    )
}
