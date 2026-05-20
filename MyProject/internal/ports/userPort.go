package ports

import (
	"MyProject/internal/entities"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// UserPort интерфейс, описывающий контракт для получения курсов криптовалют.
type UserPort interface {
	GetMax(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMin(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAvg(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetActual(ctx context.Context, titles []string) ([]entities.Coin, error)
}

// Server структура для хранения всех зависимостей HTTP сервера.
type Server struct {
	router *chi.Mux     // маршрутизатор для обработки HTTP запросов
	http   *http.Server // HTTP сервер стандартной библиотеки Go
	port   string       // порт для прослушивания ("8080")
	up     UserPort     // сервис для получения данных (бизнес-логика)
}

// NewServer создаёт и инициализирует новый HTTP сервер.
func NewServer(port string, up UserPort) *Server {
	// Создаём новый маршрутизатор chi
	r := chi.NewRouter()
	// Регистрируем middleware
	r.Use(middleware.Logger)    // логирует HTTP метод, статус, время выполнения
	r.Use(middleware.Recoverer) // ловит паники и предотвращает падение сервера
	r.Use(middleware.RequestID) // добавляет уникальный ID каждому запросу для трейсинга
	r.Use(middleware.RealIP)    // получает реальный IP клиента (при работе за прокси)

	return &Server{
		router: r,
		port:   port,
		up:     up,
	}
}

// setupRoutes настраивает все HTTP маршруты сервера.
// Должна вызваться перед Start() для регистрации обработчиков.
func (s *Server) setupRoutes() {
	// GET /health — проверка сервера
	s.router.Get("/health", s.healthHandler)

	// Группа маршрутов /coins для всех операций с курсами валют
	s.router.Route("/coins", func(r chi.Router) {
		r.Get("/max", s.getMaxHandler)
		r.Get("/min", s.getMinHandler)
		r.Get("/avg", s.getAvgHandler)
		r.Get("/actual", s.getActualHandler)
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) getMaxHandler(w http.ResponseWriter, r *http.Request) {
	s.handleCoinRequest(w, r, func(ctx context.Context, titles []string) (interface{}, error) {
		return s.up.GetMax(ctx, titles)
	})
}

func (s *Server) getMinHandler(w http.ResponseWriter, r *http.Request) {
	s.handleCoinRequest(w, r, func(ctx context.Context, titles []string) (interface{}, error) {
		return s.up.GetMin(ctx, titles)
	})
}

func (s *Server) getAvgHandler(w http.ResponseWriter, r *http.Request) {
	s.handleCoinRequest(w, r, func(ctx context.Context, titles []string) (interface{}, error) {
		return s.up.GetAvg(ctx, titles)
	})
}

func (s *Server) getActualHandler(w http.ResponseWriter, r *http.Request) {
	s.handleCoinRequest(w, r, func(ctx context.Context, titles []string) (interface{}, error) {
		return s.up.GetActual(ctx, titles)
	})
}

// coinHandlerFunc  для обработки запроса к данным о валютах.

type coinHandlerFunc func(ctx context.Context, titles []string) (interface{}, error)

// handleCoinRequest универсальный обработчик для всех запросов о курсах валют.
func (s *Server) handleCoinRequest(w http.ResponseWriter, r *http.Request, fn coinHandlerFunc) {
	// Получаем параметр "titles" из URL (например, ?titles=BTC,ETH)
	q := r.URL.Query().Get("titles")
	if q == "" {

		http.Error(w, "отсутствует параметр titles", http.StatusBadRequest)
		return
	}

	titles := strings.Split(q, ",")

	res, err := fn(r.Context(), titles)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {

		http.Error(w, "ошибка кодирования ответа", http.StatusInternalServerError)
	}
}

func (s *Server) Start() error {

	s.setupRoutes()

	s.http = &http.Server{
		Addr:           ":" + s.port, // адрес и порт для прослушивания
		Handler:        s.router,     // наш маршрутизатор
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Сервер запускается на порту %s", s.port)

	// горутина чтобы не блокировать текущий поток

	go func() {
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Останавливаем сервер...")
	return s.http.Shutdown(ctx)
}

// NewUserRouter создаёт новый маршрутизатор chi для совместимости с наследованным кодом.

func NewUserRouter(up UserPort) http.Handler {
	server := NewServer("8080", up)
	server.setupRoutes()
	return server.router
}
