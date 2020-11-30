package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var (
	// errNilDriver is returned by call any method which use connection to database if connection is nil
	errNilDriver = errors.New("nil pg driver")
)

type Client interface {
	Init(db *sql.DB) error
	GetClient() (db *sql.DB, err error)
	Close()
	Ping(ctx context.Context) error
}

type client struct {
	db      *sql.DB
	timeout time.Duration
}

func (c *client) GetClient() (db *sql.DB, err error) {
	if c.db == nil {
		err = errNilDriver
		return
	}
	db = c.db
	return
}

func (c *client) Init(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}
	c.db = db
	c.db.SetConnMaxLifetime(time.Minute * 5)
	c.db.SetMaxIdleConns(0)
	c.db.SetMaxOpenConns(10)
	return nil
}

func (c *client) Close() {
	err := c.db.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *client) Ping(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	err = c.db.PingContext(ctx)
	return
}

func NewClient(timeout time.Duration) Client {
	return &client{
		timeout: timeout,
	}
}
