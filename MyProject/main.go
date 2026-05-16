package main

import (
	"MyProject/internal/adapters"
	"MyProject/internal/cases"
	"MyProject/internal/ports"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/myproject?sslmode=disable"
	}

	storage, err := adapters.NewPostgresStorage(ctx, dbURL)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer storage.Close()

	client, err := adapters.NewClient()
	if err != nil {
		log.Fatalf("Ошибка создания клиента: %v", err)
	}

	// Инициализация сервиса
	service, err := cases.NewService(client, storage)
	if err != nil {
		log.Fatalf("Ошибка создания сервиса: %v", err)
	}

	fmt.Printf("Сервис готов к работе: %+v\n", service)

	// Инициализация HTTP роутера и запуск сервера
	router := ports.NewUserRouter(service)
	srv := &http.Server{Addr: ":8080", Handler: router}
	go func() {
		log.Printf("HTTP сервер слушает на %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ошибка HTTP сервера: %v", err)
		}
	}()

	defer func() {
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxShutdown); err != nil {
			log.Printf("ошибка завершения HTTP сервера: %v", err)
		}
	}()

	runFetch(service, 5*time.Minute)
}

func runFetch(service *cases.Service, interval time.Duration) {
	ctx := context.Background()

	if err := service.FetchRates(ctx); err != nil {
		log.Printf("FetchRates при старте завершился с ошибкой: %v", err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := service.FetchRates(ctx); err != nil {
			log.Printf("Ошибка фонового обновления курсов: %v", err)
		}
	}
}
