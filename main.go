package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator"

	"go-food-store/config"
	"go-food-store/helpers"
)

func main() {
	db := config.ConnectToDatabase()
	validate := validator.New()

	config.MigrateAllTable(db)
	config.SeedAdmin(db)

	server := &http.Server{
		Addr:              config.ServerAddr(),
		Handler:           config.SetupModel(db, validate),
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
	}

	// SIGTERM adalah yang dikirim Docker dan Kubernetes saat menghentikan
	// container. Tanpa penanganan ini proses mati seketika di tengah request
	// yang sedang berjalan — termasuk di tengah transaksi database.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server berjalan", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			helpers.PanicHelper(err)
		}
	}()

	<-ctx.Done()
	slog.Info("sinyal berhenti diterima, menunggu request selesai")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown tidak selesai tepat waktu", "error", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
	slog.Info("server berhenti")
}
