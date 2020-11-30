package server

import (
	"context"
	"github.com/0ndreu/test_kamaz/internal/models"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
)

type service interface {
	IsAlive(ctx context.Context) error
	CreateMultipleGoods(ctx context.Context, good []models.Good) error
}

// Server interface
type Server interface {
	CheckLiveness(w http.ResponseWriter, r *http.Request)
}

type server struct {
	ctx context.Context
	svc service
	log zerolog.Logger
}

func NewServer(ctx context.Context, svc service, log zerolog.Logger) http.Handler {
	srv := &server{
		ctx: ctx,
		svc: svc,
		log: log,
	}

	r := mux.NewRouter()
	sv1 := r.PathPrefix("/api/v1").Subrouter()
	r.HandleFunc("/health", srv.CheckLiveness).Methods(http.MethodGet)
	sv1.HandleFunc("/goods/add", srv.CreateMultipleGoods).Methods(http.MethodPost)

	return r
}
