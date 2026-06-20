package http

import (
	"MyProject/pkg/dto"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server структура для хранения всех зависимостей HTTP сервера.
type Server struct {
	router  *chi.Mux     // маршрутизатор для обработки HTTP запросов
	http    *http.Server // HTTP сервер стандартной библиотеки Go
	port    string       // порт для прослушивания ("8080")
	service Service      // сервис для получения данных (бизнес-логика)
}

// NewServer создаёт и инициализирует новый HTTP сервер.
func NewServer(port string, service Service) (*Server, error) {
	if port == "" {
		return nil, errors.New("порт не указан")
	}
	if service == nil {
		return nil, errors.New("сервис не указан")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.NotFound(http.NotFound)

	return &Server{
		router:  r,
		port:    port,
		service: service,
	}, nil
}

func (s *Server) StartServer() error {

	s.router.Get("v1.0.0/coins/max", s.GetMax)
	s.router.Get("v1.0.0/coins/min", s.GetMin)
	s.router.Get("v1.0.0/coins/avg", s.GetAvg)
	s.router.Get("v1.0.0/coins/actual", s.GetActual)

	// Запускаем сервер
	s.http = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	log.Printf("Сервер запущен на порту %s", s.port)
	return s.http.ListenAndServe()
}

func (s *Server) GetMax(w http.ResponseWriter, r *http.Request) {
	titles, err := getTitles(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := s.service.GetMax(r.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := make(dto.CoinsDTO, len(result))
	for i, coin := range result {
		response[i] = dto.CoinDTO{
			Title: coin.Title,
			Rate:  coin.Rate,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

func (s *Server) GetMin(w http.ResponseWriter, r *http.Request) {
	titles, err := getTitles(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := s.service.GetMin(r.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := make(dto.CoinsDTO, len(result))
	for i, coin := range result {
		response[i] = dto.CoinDTO{
			Title: coin.Title,
			Rate:  coin.Rate,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

func (s *Server) GetAvg(w http.ResponseWriter, r *http.Request) {
	titles, err := getTitles(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := s.service.GetAvg(r.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	response := make(dto.CoinsDTO, len(result))
	for i, coin := range result {
		response[i] = dto.CoinDTO{
			Title: coin.Title,
			Rate:  coin.Rate,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

func (s *Server) GetActual(w http.ResponseWriter, r *http.Request) {
	titles, err := getTitles(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := s.service.GetActual(r.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := make(dto.CoinsDTO, len(result))
	for i, coin := range result {
		response[i] = dto.CoinDTO{
			Title: coin.Title,
			Rate:  coin.Rate,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}
func getTitles(r *http.Request) ([]string, error) {
	titles := r.URL.Query().Get("titles")

	if titles == "" {
		return nil, errors.New("Не введена валюта")
	}
	return strings.Split(titles, ","), nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.http.Shutdown(ctx)
}
