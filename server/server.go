package server

import (
	"net/http"

	"github.com/go-kit/log"
	h "github.com/go-url-shortener/api"
	"github.com/go-url-shortener/repository/cachedpostgres"
	"github.com/go-url-shortener/repository/postgres"
	"github.com/go-url-shortener/repository/redis"
	"github.com/go-url-shortener/shortener"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewServer(logger log.Logger) http.Handler {
	r := mux.NewRouter()
	// for now, just postgres
	repo, err := postgres.NewPostgresRepository("postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable")
	redisRepo, err := redis.NewRedisRepository("redis://127.0.0.1:6379")
	if err != nil {
		logger.Log("error", err)
	}

	cachedRepository := cachedpostgres.NewCachedRepository(repo, redisRepo)

	service := shortener.NewRedirectService(cachedRepository)
	service = NewLoggingMiddleware(logger, service)
	handler := h.NewHandler(service)

	r.Use(TrackApiCalls)
	r.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	r.HandleFunc("/{code}", handler.Get).Methods("GET")
	r.HandleFunc("/store", handler.Post).Methods("POST")

	return r
}
