package main

import (
	"MyProject/internal/adapters"
	"MyProject/internal/cases"
	"context"
	"fmt"
	"log"
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
