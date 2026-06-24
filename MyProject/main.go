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

// @title           MyProject Crypto API
// @version         1.0
// @description     API for retrieving cryptocurrency exchange rates and statistics on average, minimum, and maximum values.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// GetMax godoc
// @Summary      Fetch maximum currency rates
// @Description  Returns the maximum rates for the specified currencies over all time
// @Tags         coins
// @Accept       json
// @Produce      json
// @Param        titles path string true "List of currencies separated by commas (e.g., BTC,ETH,USD)"
// @Success      200 {array} dto.CoinDTO "Successful response with maximum rates"
// @Failure      400 {string} string "Invalid request (missing titles parameter)"
// @Failure      404 {string} string "No coins found"
// @Failure      500 {string} string "Internal server error"
// @Router       /coins/get_max/{titles} [get]

// GetMin godoc
// @Summary      Fetch minimum currency rates
// @Description  Returns the minimum rates for the specified currencies over all time
// @Tags         coins
// @Accept       json
// @Produce      json
// @Param        titles path string true "List of currencies separated by commas (e.g., BTC,ETH,USD)"
// @Success      200 {array} dto.CoinDTO "Successful response with minimum rates"
// @Failure      400 {string} string "Invalid request (missing titles parameter)"
// @Failure      404 {string} string "No coins found"
// @Failure      500 {string} string "Internal server error"
// @Router       /coins/get_min/{titles} [get]

// GetAvg godoc
// @Summary      Fetch average currency rates
// @Description  Returns the average rates for the specified currencies over all time
// @Tags         coins
// @Accept       json
// @Produce      json
// @Param        titles path string true "List of currencies separated by commas (e.g., BTC,ETH,USD)"
// @Success      200 {array} dto.CoinDTO "Successful response with average rates"
// @Failure      400 {string} string "Invalid request (missing titles parameter)"
// @Failure      404 {string} string "No coins found"
// @Failure      500 {string} string "Internal server error"
// @Router       /coins/get_avg/{titles} [get]

// GetActual godoc
// @Summary      Fetch actual currency rates
// @Description  Returns the latest actual rates for the specified currencies
// @Tags         coins
// @Accept       json
// @Produce      json
// @Param        titles path string true "List of currencies separated by commas (e.g., BTC,ETH,USD)"
// @Success      200 {array} dto.CoinDTO "Successful response with actual rates"
// @Failure      400 {string} string "Invalid request (missing titles parameter)"
// @Failure      404 {string} string "No coins found"
// @Failure      500 {string} string "Internal server error"
// @Router       /coins/get_actual/{titles} [get]
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
