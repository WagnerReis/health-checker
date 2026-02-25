package main

import (
	"context"
	"errors"
	"fmt"
	"health-checker/config"
	"health-checker/internal/application/usecases"
	"health-checker/internal/infra/cryptography"
	router "health-checker/internal/infra/http"
	"health-checker/internal/infra/http/handlers"
	"health-checker/internal/infra/logger"
	dbutils "health-checker/internal/infra/persistence/database"
	"health-checker/internal/infra/persistence/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger, err := logger.NewZapLogger()
	if err != nil {
		panic(err)
	}

	cfg := config.LoadConfig()

	tokenGenerator := cryptography.NewJWTTokenGenerator()
	hasher := cryptography.NewBcrypterHasher()
	sha256Hash := cryptography.NewSHA256Hash()

	db, err := dbutils.NewPool(context.Background())
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error creating database pool: %v", err))
	}

	// Repositories
	userRepository := postgres.NewUserRepository(db)
	refreshTokenRepository := postgres.NewRefreshTokenRepository(db)
	monitorRepository := postgres.NewMonitorRepository(db)

	// UseCases
	// Auth
	signUpUseCase := usecases.NewSignUpUseCase(userRepository, refreshTokenRepository, hasher, tokenGenerator, sha256Hash, *cfg, logger)
	loginUseCase := usecases.NewLoginUseCase(userRepository, refreshTokenRepository, hasher, tokenGenerator, sha256Hash, *cfg, logger)
	logoutUseCase := usecases.NewLogoutUseCase(userRepository, refreshTokenRepository, sha256Hash, logger)
	refreshUseCase := usecases.NewRefreshUseCase(userRepository, refreshTokenRepository, tokenGenerator, sha256Hash, *cfg, logger)

	// Monitor
	createMonitorUseCase := usecases.NewCreateMonitorUseCase(monitorRepository, logger)
	getMonitorsUseCase := usecases.NewGetMonitorsUseCase(monitorRepository, logger)

	// Handlers
	authHandler := handlers.NewAuthHandler(*signUpUseCase, *loginUseCase, *logoutUseCase, *refreshUseCase)
	monitorHandler := handlers.NewMonitorHandler(*createMonitorUseCase, *getMonitorsUseCase)

	// Router
	appRouter := router.NewAppRouter(authHandler, monitorHandler)
	router := appRouter.InitializeRoutes()

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Server is running on port: " + cfg.Port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(fmt.Sprintf("Error starting server: %v", err))
		}
		logger.Info("Stopped serving new connections")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, shutdownRelease := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownRelease()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	logger.Info("Gracefull shutdown complete.")
}
