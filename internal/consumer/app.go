package consumer

import (
	"context"
	"log"
	"message-service/internal/messages"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	DB    *pgxpool.Pool
	Kafka *kafka.Reader
}

type App struct {
	config  *Config
	service *Service
}

func NewApp(config *Config) *App {
	return &App{
		config:  config,
		service: NewService(messages.NewService(config.DB)),
	}
}

func (app *App) Run() error {
	reader := app.config.Kafka

	for {
		m, err := reader.ReadMessage(context.Background())

		if err != nil {
			break
		}

		go app.service.Handle(context.Background(), m)

		log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	return reader.Close()
}
