package ports

import (
	"MyProject/internal/entities"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type UserPort interface {
	GetMax(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMin(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAvg(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetPercent(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetLast(ctx context.Context, titles []string) ([]entities.Coin, error)
}

func NewUserRouter(up UserPort) http.Handler {
	r := chi.NewRouter()

	r.Get("/max", func(w http.ResponseWriter, req *http.Request) {
		handleUserRequest(w, req, func(ctx context.Context, titles []string) (interface{}, error) {
			return up.GetMax(ctx, titles)
		})
	})

	r.Get("/min", func(w http.ResponseWriter, req *http.Request) {
		handleUserRequest(w, req, func(ctx context.Context, titles []string) (interface{}, error) {
			return up.GetMin(ctx, titles)
		})
	})

	r.Get("/avg", func(w http.ResponseWriter, req *http.Request) {
		handleUserRequest(w, req, func(ctx context.Context, titles []string) (interface{}, error) {
			return up.GetAvg(ctx, titles)
		})
	})

	r.Get("/perc", func(w http.ResponseWriter, req *http.Request) {
		handleUserRequest(w, req, func(ctx context.Context, titles []string) (interface{}, error) {
			return up.GetPercent(ctx, titles)
		})
	})

	r.Get("/last", func(w http.ResponseWriter, req *http.Request) {
		handleUserRequest(w, req, func(ctx context.Context, titles []string) (interface{}, error) {
			return up.GetLast(ctx, titles)
		})
	})

	return r
}

type userHandlerFunc func(ctx context.Context, titles []string) (interface{}, error)

func handleUserRequest(w http.ResponseWriter, req *http.Request, fn userHandlerFunc) {
	q := req.URL.Query().Get("titles")
	if q == "" {
		http.Error(w, "отсутствует параметр titles", http.StatusBadRequest)
		return
	}
	titles := strings.Split(q, ",")

	res, err := fn(req.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "ошибка кодирования ответа", http.StatusInternalServerError)
	}
}
