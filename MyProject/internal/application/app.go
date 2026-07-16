package application

import (
	"MyProject/deploy/config"
	"MyProject/internal/adapters/api_client/coindesk"
	"MyProject/internal/adapters/storage/postgres"
	"MyProject/internal/cases"
	"MyProject/internal/ports/api_user/http"
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run() {
	ctx := context.Background()

	// Получаем порт из переменной окружения PORT или используем 8080 по умолчанию

	port := a.cfg.Port
	if port == "" {
		port = "8080"
	}

	// Получаем строку подключения к БД из переменной окружения или используем локальный PostgreSQL
	dbURL := a.cfg.URL
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/myproject?sslmode=disable"
	}

	storage, err := postgres.NewPostgresStorage(ctx, dbURL)
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}

	defer storage.Close()

	// Создаём HTTP клиент
	client, err := coindesk.NewClient(
		a.cfg.ExternalAPICoindeskBaseURL,
		a.cfg.ExternalAPICoindeskTimeout,
		a.cfg.ExternalAPICoindeskCurrency,
	)
	if err != nil {
		log.Fatalf("Ошибка создания клиента: %v", err)
	}

	// Инициализация сервиса (бизнес-логика)
	service, err := cases.NewService(client, storage)
	if err != nil {
		log.Fatalf("Ошибка создания сервиса: %v", err)
	}

	// Service реализует интерфейс HTTP-порта сервера — передаём его в HTTP-слой

	go startCron(service, a.cfg.UpdateInterval)

	server, err := http.NewServer(port, service)
	if err != nil {
		log.Fatalf("Ошибка создания HTTP сервера: %v", err)
	}

	if err := server.StartServer(); err != nil {
		log.Fatal(err)
	}

}

// runFetch запускает периодическое обновление курсов валют из внешнего API.

func runFetch(service *cases.Service) {
	ctx := context.Background()

	if err := service.FetchRates(ctx); err != nil {
		log.Printf("FetchRates при старте завершился с ошибкой: %v", err)
	}
}

func startCron(service *cases.Service, interval string) {

	c := cron.New()

	c.AddFunc("@every "+interval, func() {
		runFetch(service)
	})
	c.Start()

	log.Println("Cron задача запущена (каждые " + interval + " секунд)")
}
