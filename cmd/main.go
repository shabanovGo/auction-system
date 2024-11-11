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

    createUserUC := userUseCase.NewCreateUserUseCase(userRepo)
    getUserUC := userUseCase.NewGetUserUseCase(userRepo)
    updateUserUC := userUseCase.NewUpdateUserUseCase(userRepo)
    deleteUserUC := userUseCase.NewDeleteUserUseCase(userRepo)
    getAllUsersUC := userUseCase.NewGetAllUserUseCase(userRepo)
    updateBalanceUC := userUseCase.NewUpdateBalanceUseCase(userRepo)

    application := app.NewApp(
        cfg,
        createUserUC,
        getUserUC,
        updateUserUC,
        deleteUserUC,
        getAllUsersUC,
        updateBalanceUC,
    )

    log.Println("Starting application...")
    if err := application.Run(ctx); err != nil {
        log.Fatalf("Error running app: %v", err)
    }
}
