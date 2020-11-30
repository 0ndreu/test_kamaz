package service

import (
	"context"
	"github.com/0ndreu/test_kamaz/internal/models"
)

type storage interface {
	CreateMultipleGoods(ctx context.Context, good []models.Good) error
}

type client interface {
	Ping(ctx context.Context) (err error)
}

type Service interface {
	IsAlive(ctx context.Context) error
	CreateMultipleGoods(ctx context.Context, good []models.Good) error
}

type service struct {
	storage storage
	client  client
}

func (s *service) CreateMultipleGoods(ctx context.Context, good []models.Good) error {
	return s.storage.CreateMultipleGoods(ctx, good)
}

func (s *service) IsAlive(ctx context.Context) error {
	return s.client.Ping(ctx)
}

func NewService(storage storage, client client) Service {
	return &service{
		storage: storage,
		client:  client,
	}
}
