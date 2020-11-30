package storage

import (
	"context"
	"github.com/0ndreu/test_kamaz/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"time"
)

type Storage interface {
	CreateMultipleGoods(ctx context.Context, good []models.Good) error
}

type storage struct {
	client  Client
	timeout time.Duration
}

func (s *storage) PostNewGood(ctx context.Context, good []models.Good) error {
	conn, err := s.client.GetClient()
	if err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, "", pq.Array(good))
	return err
}

func (s *storage) CreateMultipleGoods(ctx context.Context, goods []models.Good) error {
	conn, err := s.client.GetClient()
	if err != nil {
		return err
	}

	query := sq.Insert("public.good").
		Columns("amount", "method", "object", "price", "quantity")
	for _, i := range goods {
		query = query.Values(
			i.Amount,
			i.Method,
			i.Object,
			i.Price,
			i.Quantity)
	}

	queryInsert, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, queryInsert, args...)
	return err
}

func NewStorage(client Client, timeout time.Duration) Storage {
	return &storage{
		client:  client,
		timeout: timeout,
	}
}
