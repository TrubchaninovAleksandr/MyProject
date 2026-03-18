package main

import (
	"MyProject/internal/adapters"
	"MyProject/internal/cases"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	storage := adapters.NewPostgresStorage(db)
	client := adapters.NewClient()

	// Инициализация сервиса
	service, err := cases.NewService(client, storage)
	if err != nil {
		log.Fatalf("Ошибка создания сервиса: %v", err)
	}

	fmt.Printf("Сервис готов к работе: %+v\n", service)

	// Запускаем фоновое обновление курсов каждые 5 минут
	runFetch(service, 5*time.Minute)
}

func runFetch(service *cases.Service, interval time.Duration) {
	ctx := context.Background()

	// Первый fetch выполняем сразу при старте
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
