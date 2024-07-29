package messages

import (
	"encoding/base64"
	"fmt"
	"message-service/internal/statistics"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

type Auth struct {
	Username string
	Password string
}

type Config struct {
	Port string
	Host string
	Auth
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

func BasicAuth(next http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		authHeader = strings.TrimPrefix(authHeader, "Basic ")

		auth, err := base64.StdEncoding.DecodeString(authHeader)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		authUsername, authPassword, found := strings.Cut(string(auth), ":")

		if !found || authUsername != username || authPassword != password {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) Run() error {
	messagesHandler := NewHandler(app.config.DB, app.config.Kafka)
	statisticsHandler := statistics.NewHandler(app.config.DB)

	mux := http.NewServeMux()

	mux.Handle("GET /messages/{id}", http.HandlerFunc(messagesHandler.GetMessage))
	mux.Handle("GET /messages/{id}/", http.HandlerFunc(messagesHandler.GetMessage))

	mux.Handle("POST /messages", http.HandlerFunc(messagesHandler.PostMessage))
	mux.Handle("POST /messages/", http.HandlerFunc(messagesHandler.PostMessage))

	mux.Handle("GET /stats", http.HandlerFunc(statisticsHandler.GetStatistics))

	mux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	authenticatedMux := BasicAuth(mux, app.config.Auth.Username, app.config.Auth.Password)

	return http.ListenAndServe(fmt.Sprintf("%s:%s", app.config.Host, app.config.Port), authenticatedMux)
}
