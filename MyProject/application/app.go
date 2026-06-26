package app

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

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
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

	// Попытаемся подключиться к Postgres. Если подключение не удалось — завершаем программу с ошибкой.
	var storage cases.Storage
	pgStorage, err := postgres.NewPostgresStorage(ctx, dbURL)
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}
	// Если подключение успешно — используем Postgres и закроем пул при выходе.
	defer pgStorage.Close()
	storage = pgStorage

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

	// Service реализует интерфейс HTTP-порта сервера — передаём его в HTTP-слой как контракт.
	var userPort http.Service = service

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
