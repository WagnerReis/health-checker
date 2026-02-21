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

	db, err := dbutils.NewPool(context.Background())
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error creating database pool: %v", err))
	}

	// Repositories
	userRepository := postgres.NewUserRepository(db)

	// UseCases
	signUpUseCase := usecases.NewSignUpUseCase(userRepository, hasher, tokenGenerator, *cfg, logger)
	loginUseCase := usecases.NewLoginUseCase(userRepository, hasher, tokenGenerator, *cfg, logger)

	// Handlers
	authHandler := handlers.NewAuthHandler(*signUpUseCase, *loginUseCase)

	// Router
	appRouter := router.NewAppRouter(authHandler)
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
