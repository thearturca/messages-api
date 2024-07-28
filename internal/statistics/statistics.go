package statistics

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	service   *Service
	validator *validator.Validate
	db        *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{
		service:   NewService(db),
		validator: validator.New(),
		db:        db,
	}
}

func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	fromAsString := r.URL.Query().Get("from")
	toAsString := r.URL.Query().Get("to")

	var from *time.Time
	var to *time.Time

	if fromAsString != "" {
		fromParsed, err := time.Parse(time.RFC3339, fromAsString)

		if err != nil {
			http.Error(w, "invalid from date", http.StatusBadRequest)
			return
		}

		from = &fromParsed
	}

	if toAsString != "" {
		toParsed, err := time.Parse(time.RFC3339, toAsString)

		if err != nil {
			http.Error(w, "invalid to date", http.StatusBadRequest)
			return
		}

		to = &toParsed
	}

	statistics, err := h.service.GetStatistics(r.Context(), from, to)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	marshalledStatistics, err := json.Marshal(statistics)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshalledStatistics)
}
