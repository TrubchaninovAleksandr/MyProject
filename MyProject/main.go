package main

import (
	"MyProject/internal/adapters/api_client/coindesk"
	"MyProject/internal/adapters/storage/postgres"
	"MyProject/internal/cases"
	"MyProject/internal/ports/api_user/http"
	"context"
	"log"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	// Получаем порт из переменной окружения PORT или используем 8080 по умолчанию
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Получаем строку подключения к БД из переменной окружения или используем локальный PostgreSQL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/myproject?sslmode=disable"
	}

	storage, err := postgres.NewPostgresStorage(ctx, dbURL)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer storage.Close()

	// Создаём HTTP клиент для получения курсов от внешнего API
	client, err := coindesk.NewClient()
	if err != nil {
		log.Fatalf("Ошибка создания клиента: %v", err)
	}

	// Инициализация сервиса (бизнес-логика)
	service, err := cases.NewService(client, storage)
	if err != nil {
		log.Fatalf("Ошибка создания сервиса: %v", err)
	}

	// Service реализует cases.UserPort — передаём его в HTTP-слой как контракт.
	var userPort cases.UserPort = service

	go runFetch(service, 5*time.Minute)

	server, err := http.NewServer(port, userPort)
	if err != nil {
		log.Fatalf("Ошибка создания HTTP сервера: %v", err)
	}

	if err := server.StartServer(); err != nil {
		log.Fatal(err)
	}
}

// runFetch запускает периодическое обновление курсов валют из внешнего API.
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
