package main

import (
	"context"
	"errors"
	"fmt"
	"health-checker/config"
	router "health-checker/internal/infra/http"
	"health-checker/internal/infra/logger"
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
	appRouter := router.NewAppRouter()
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
