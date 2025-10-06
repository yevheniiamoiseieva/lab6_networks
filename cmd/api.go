package main

import (
	"context"
	"fmt"
	"laba6/internal/app"
	"laba6/pkg/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfiguration()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, err := app.NewApplication(ctx, cfg)
	if err != nil {
		panic("Failed to start application: " + err.Error())
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server starting on port %s...\n", cfg.Application.Port)
		if err := application.Start(); err != nil && err != http.ErrServerClosed {
			panic("Server start error: " + err.Error())
		}
	}()

	<-stop
	fmt.Println("Received shutdown signal, shutting down...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := application.Shutdown(ctxShutdown); err != nil {
		panic("Error during shutdown: " + err.Error())
	}

	fmt.Println("Server shutdown completed")
}
