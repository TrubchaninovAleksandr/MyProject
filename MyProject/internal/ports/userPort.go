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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	return &Server{
		router: r,
		port:   port,
		up:     up,
	}
}
func (s *Server) Start() error {

	s.router.Get("/coins/max", s.GetMax)
	s.router.Get("/coins/min", s.GetMin)
	s.router.Get("/coins/avg", s.GetAvg)
	s.router.Get("/coins/actual", s.GetActual)

	// Запускаем сервер
	s.http = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	log.Printf("Сервер запущен на порту %s", s.port)
	return s.http.ListenAndServe()
}

func (s *Server) GetMax(w http.ResponseWriter, r *http.Request) {
	titles := getTitles(r)
	result, err := s.up.GetMax(r.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (s *Server) GetMin(w http.ResponseWriter, r *http.Request) {
	titles := getTitles(r)
	result, err := s.up.GetMin(r.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (s *Server) GetAvg(w http.ResponseWriter, r *http.Request) {
	titles := getTitles(r)
	result, err := s.up.GetAvg(r.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (s *Server) GetActual(w http.ResponseWriter, r *http.Request) {
	titles := getTitles(r)
	result, err := s.up.GetActual(r.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}
func getTitles(r *http.Request) []string {
	titles := r.URL.Query().Get("titles")
	if titles == "" {
		return []string{}
	}
	return strings.Split(titles, ",")
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.http.Shutdown(ctx)
}
