package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	core_api "github.com/emp2ty0/coal-mine/internal/core/api"
	"github.com/emp2ty0/coal-mine/internal/core/domain"
	"github.com/gin-gonic/gin"
)

func main() {
	mtx := sync.Mutex{}
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	enterprise := domain.NewEnterprise(ctx, &mtx)
	go enterprise.Run()
	httpHandlers := core_api.NewHTTPHandlers(enterprise, logger, ctx)

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.GET("/miners/wages", httpHandlers.MinersWages)
		api.GET("/miners/hire", httpHandlers.MinerHire)

		api.GET("/enterprise/coal", httpHandlers.GetCoal)
		api.GET("/enterprise/miners", httpHandlers.MinersShow)
		api.GET("/enterprise/devices", httpHandlers.ShowDevices)

		api.GET("/device/info", httpHandlers.DeviceInfo)
		api.GET("/device/buy", httpHandlers.BuyDevice)

		api.GET("/cancel", httpHandlers.CancelContext)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	chanError := make(chan error, 1)

	go func() {
		defer close(chanError)
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			chanError <- err
		}
	}()

	select {
	case err := <-chanError:
		if err != nil {
			logger.Warn("server error" + err.Error())
			os.Exit(1)
		}
	case <-ctx.Done():
		logger.Warn("shutdown HTTP server")

		shutDownCtx, cancel := context.WithTimeout(context.Background(), 30)

		defer cancel()

		if err := server.Shutdown(shutDownCtx); err != nil {
			logger.Error("shutdown error" + err.Error())
			_ = server.Close()
		}

		logger.Info("HTTP server stopped")
	}
}
