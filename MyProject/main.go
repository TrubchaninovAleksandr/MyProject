package main

import (
	"MyProject/internal/adapters"
	"MyProject/internal/cases"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func main() {

	var db *sql.DB
	storage := adapters.NewPostgresStorage(db)
	client := adapters.NewClient()

	// Инициализация сервиса
	service, err := cases.NewService(client, storage)
	if err != nil {
		log.Fatalf("Ошибка создания сервиса: %v", err)
	}

	fmt.Printf("Сервис готов к работе: %+v\n", service)

	// Запускаем фоновое обновление курсов
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
