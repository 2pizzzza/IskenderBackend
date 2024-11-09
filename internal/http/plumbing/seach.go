package plumbing

import (
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
	const op = "handler.Search"
	log := s.log.With(
		slog.String("op", op),
	)

	isProducer := r.URL.Query().Get("is_producer")
	searchQuery := r.URL.Query().Get("q")
	code := r.URL.Query().Get("lang")

	var isProducerVal *bool
	if isProducer != "" {
		val, err := strconv.ParseBool(isProducer)
		if err != nil {
			log.Error("Invalid isProducer value", slog.String("isProducer", isProducer), sl.Err(err))
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid isProducer value"}, http.StatusInternalServerError)
			return
		}
		isProducerVal = &val
	}

	result, err := s.service.Search(r.Context(), code, isProducerVal, searchQuery)
	if err != nil {
		log.Error("Failed to execute search", sl.Err(err))
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to execute search"}, http.StatusInternalServerError)

		return
	}

	utils.WriteResponseBody(w, result, http.StatusOK)
}
