package plumbing

import (
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/2pizzzza/plumbing/internal/utils"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
)

// Search performs a search for collections and items based on query parameters
// @Summary Search for collections and items based on language, producer status, and search query
// @Description Performs a search for collections and items based on the provided parameters such as language, producer status, and search query.
// @Tags search
// @Accept  json
// @Produce  json
// @Param  lang  query  string  false  "Language code"
// @Param  is_producer  query  bool  false  "Filter by producer status"
// @Param  is_painted  query  bool  false  "Filter by painted"
// @Param  q  query  string  false  "Search query"
// @Success 200 {object} models.PopularResponse "Search results"
// @Failure 400 {object} models.ErrorMessage "Bad request - invalid query parameters"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /search [get]
func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
	const op = "handler.Search"
	log := s.log.With(
		slog.String("op", op),
	)

	isProducer := r.URL.Query().Get("is_producer")
	isPainted := r.URL.Query().Get("is_painted")
	searchQuery := r.URL.Query().Get("q")
	code := r.URL.Query().Get("lang")

	decodedQuery, err := url.QueryUnescape(searchQuery)
	if err != nil {
		fmt.Println("Error decoding query:", err)
	}
	fmt.Println("Decoded query:", decodedQuery)

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

	var isPaintedVal *bool
	if isPainted != "" {
		val, err := strconv.ParseBool(isPainted)
		if err != nil {
			log.Error("Invalid isPainted value", slog.String("isPainted", isPainted), sl.Err(err))
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid isPainted value"}, http.StatusInternalServerError)
			return
		}
		isPaintedVal = &val
	}

	result, err := s.service.Search(r.Context(), code, isProducerVal, isPaintedVal, decodedQuery)
	if err != nil {
		if errors.Is(err, storage.ErrCollectionNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Not Found"}, http.StatusNotFound)
			return
		}
		log.Error("Failed to execute search", sl.Err(err))
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to execute search"}, http.StatusInternalServerError)

		return
	}

	utils.WriteResponseBody(w, result, http.StatusOK)
}
