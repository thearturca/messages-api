package messagesProcessor

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Kafka struct {
		Reader *kafka.Reader
		Writer *kafka.Writer
	}
}

type App struct {
	config  *Config
	service *Service
}

func NewApp(config *Config) *App {
	return &App{
		service: NewService(config.Kafka.Writer),
	}
}

func (app *App) Run() error {
	reader := app.config.Kafka.Reader

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
