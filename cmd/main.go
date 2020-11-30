package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/0ndreu/test_kamaz/internal/server"
	"github.com/0ndreu/test_kamaz/internal/service"
	"github.com/0ndreu/test_kamaz/internal/storage"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	// Server
	Port     int           `envconfig:"PORT" required:"true"`
	RTimeout time.Duration `envconfig:"READ_TIMEOUT" required:"true"`
	WTimeout time.Duration `envconfig:"WRITE_TIMEOUT" required:"true"`

	// Database
	DBConnStr string        `envconfig:"DB_CONN_STR" required:"false"`
	DBTimeout time.Duration `envconfig:"DB_TIMEOUT" required:"true"`
}

func main() {
	log := zerolog.New(os.Stdout).With().Logger()

	var conf config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	db, err := sql.Open("postgres", conf.DBConnStr)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	client := storage.NewClient(conf.DBTimeout)
	if err := client.Init(db); err != nil {
		log.Fatal().Msg(err.Error())
	}

	storage := storage.NewStorage(client, conf.DBTimeout)
	service := service.NewService(storage, client)
	ctx, cancel := context.WithCancel(context.Background())
	r := server.NewServer(ctx, service, log)

	api := http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Port),
		ReadTimeout:  conf.RTimeout,
		Handler:      r,
		WriteTimeout: conf.WTimeout,
	}

	go func() {
		log.Info().Msgf("Listen to API on %d port", conf.Port)
		err := api.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Listen to API: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-quit
	cancel()

}
