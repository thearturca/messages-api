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
		config:  config,
		service: NewService(config.Kafka.Writer),
	}
}

func (app *App) Run() error {
	reader := app.config.Kafka.Reader

	ctx := context.Background()
	for {
		m, err := reader.ReadMessage(ctx)

		if err != nil {
			log.Println("error reading message: ", err)
			break
		}

		go app.service.Handle(ctx, m)

		log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	return reader.Close()
}
