package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nightmaker00/go-tasks-api/internal/api"
	"github.com/nightmaker00/go-tasks-api/internal/config"
	"github.com/nightmaker00/go-tasks-api/internal/repository"
	"github.com/nightmaker00/go-tasks-api/internal/service"
	"github.com/nightmaker00/go-tasks-api/pkg/db/postgres"

	_ "github.com/nightmaker00/go-tasks-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Tasks API
// @version         1.0
// @description     REST API для управления задачами (CRUDL)
// @description     Реализовано на чистом Go с использованием стандартной библиотеки
// @host            localhost:8080
// @BasePath        /
// @schemes         http

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.Open(cfg.Config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	handler := api.NewHandler(taskService)

	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	handler.RegisterRoutes(mux)
	rootHandler := api.WithCORS(mux)

	server := &http.Server{
		Addr:         cfg.Server.Address + ":" + cfg.Server.Port,
		Handler:      rootHandler,
		ReadTimeout:  time.Duration(cfg.Server.Timeouts.ReadSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Timeouts.WriteSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.Timeouts.IdleSeconds) * time.Second,
	}
	//graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
