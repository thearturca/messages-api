package messages

import (
	"fmt"
	"message-service/internal/statistics"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	Port  string
	DB    *pgxpool.Pool
	Kafka *kafka.Writer
}

type App struct {
	config *Config
}

func NewApp(config *Config) *App {
	return &App{
		config: config,
	}
}

func (app *App) Run() error {
	messagesHandler := NewHandler(app.config.DB, app.config.Kafka)
	statisticsHandler := statistics.NewHandler(app.config.DB)

	mux := http.NewServeMux()

	mux.Handle("GET /messages/{id}", http.HandlerFunc(messagesHandler.GetMessage))
	mux.Handle("GET /messages/{id}/", http.HandlerFunc(messagesHandler.GetMessage))

	mux.Handle("POST /messages", http.HandlerFunc(messagesHandler.PostMessage))
	mux.Handle("POST /messages/", http.HandlerFunc(messagesHandler.PostMessage))

	mux.Handle("GET /statistics", http.HandlerFunc(statisticsHandler.GetStatistics))

	return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", app.config.Port), mux)
}
