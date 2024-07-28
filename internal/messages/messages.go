package messages

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	"message-service/internal/kafka"
)

type PostMessageDto struct {
	Text string `json:"text" validate:"required"`
}

type Handler struct {
	db            *pgxpool.Pool
	validator     *validator.Validate
	service       *Service
	writerService *kafkaService.WriterService
}

func NewHandler(db *pgxpool.Pool, kafka *kafka.Writer) *Handler {
	return &Handler{
		db:            db,
		validator:     validator.New(),
		service:       NewService(db),
		writerService: kafkaService.NewWriterService(kafka),
	}
}

func (h *Handler) GetMessage(w http.ResponseWriter, r *http.Request) {
	messageId := r.PathValue("id")

	err := h.validator.Var(messageId, "required,uuid")

	if err != nil {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	message, err := h.service.GetMessage(r.Context(), messageId)

	if err != nil {
		switch {
		case err == pgx.ErrNoRows:
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	marshalledMessage, err := json.Marshal(message)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshalledMessage)
}

func (h *Handler) PostMessage(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var dto PostMessageDto
	json.NewDecoder(r.Body).Decode(&dto)
	defer r.Body.Close()

	if err := h.validator.Struct(dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := h.service.PostMessage(r.Context(), dto.Text)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	marshalledMessage, err := json.Marshal(message)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.writerService.SendMessage(r.Context(), message)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(marshalledMessage)
}
