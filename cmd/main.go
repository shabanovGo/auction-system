package main

import (
    "context"
    "database/sql"
    "log"
    _ "github.com/lib/pq"

    "auction-system/internal/app"
    "auction-system/internal/config"
    "auction-system/internal/infrastructure/persistence/postgres"
    userUseCase "auction-system/internal/application/usecase/user"
    lotUseCase "auction-system/internal/application/usecase/lot"
	auctionUseCase "auction-system/internal/application/usecase/auction"
    handler "auction-system/internal/interfaces/grpc/handler"
)

func main() {
    ctx := context.Background()

    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    db, err := sql.Open("postgres", cfg.Database.GetDSN())
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    userRepo := postgres.NewUserRepository(db)
    lotRepo := postgres.NewLotRepository(db)
    auctionRepo := postgres.NewAuctionRepository(db)

    createUserUC := userUseCase.NewCreateUserUseCase(userRepo)
    getUserUC := userUseCase.NewGetUserUseCase(userRepo)
    updateUserUC := userUseCase.NewUpdateUserUseCase(userRepo)
    deleteUserUC := userUseCase.NewDeleteUserUseCase(userRepo)
    getAllUsersUC := userUseCase.NewGetAllUserUseCase(userRepo)
    updateBalanceUC := userUseCase.NewUpdateBalanceUseCase(userRepo)

    createLotUC := lotUseCase.NewCreateLotUseCase(lotRepo)
    getLotUC := lotUseCase.NewGetLotUseCase(lotRepo)
    updateLotUC := lotUseCase.NewUpdateLotUseCase(lotRepo)
    deleteLotUC := lotUseCase.NewDeleteLotUseCase(lotRepo)
    getLotsUC := lotUseCase.NewGetLotsUseCase(lotRepo)

	createAuctionUC := auctionUseCase.NewCreateAuctionUseCase(auctionRepo, lotRepo)
	getAuctionUC := auctionUseCase.NewGetAuctionUseCase(auctionRepo)
    updateAuctionUC := auctionUseCase.NewUpdateAuctionUseCase(auctionRepo)
    listAuctionsUC := auctionUseCase.NewListAuctionsUseCase(auctionRepo)


    userHandler := handler.NewUserHandler(
        createUserUC,
        getUserUC,
        updateUserUC,
        deleteUserUC,
        getAllUsersUC,
        updateBalanceUC,
    )

    lotHandler := handler.NewLotHandler(
        createLotUC,
        getLotUC,
        updateLotUC,
        deleteLotUC,
        getLotsUC,
    )

	auctionHandler := handler.NewAuctionHandler(
        createAuctionUC,
        getAuctionUC,
        updateAuctionUC,
        listAuctionsUC,
    )


    handlers := handler.NewHandlers(userHandler, auctionHandler, lotHandler)

    application := app.NewApp(cfg, handlers)

    log.Println("Starting application...")
    if err := application.Run(ctx); err != nil {
        log.Fatalf("Error running app: %v", err)
    }
}
