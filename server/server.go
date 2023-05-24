package server

import (
	"log"
	"net/http"

	h "github.com/go-url-shortener/api"
	"github.com/go-url-shortener/repository/cachedpostgres"
	"github.com/go-url-shortener/repository/postgres"
	"github.com/go-url-shortener/repository/redis"
	"github.com/go-url-shortener/shortener"
	"github.com/gorilla/mux"
)

func NewServer() http.Handler {
	r := mux.NewRouter()
	// for now, just postgres
	repo, err := postgres.NewPostgresRepository("postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable")
	redisRepo, err := redis.NewRedisRepository("redis://127.0.0.1:6379")
	if err != nil {
		log.Fatal(err)
	}

	cachedRepository := cachedpostgres.NewCachedRepository(repo, redisRepo)

	service := shortener.NewRedirectService(cachedRepository)
	handler := h.NewHandler(service)

	r.HandleFunc("/{code}", handler.Get).Methods("GET")
	r.HandleFunc("/store", handler.Post).Methods("POST")
	return r
}
