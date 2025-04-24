package handler

import (
	"net/http"

	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/internal/service"
)

type GoogleHandler struct {
	googleSvc service.GoogleService
}

func NewGoogleHandler(googleSvc service.GoogleService) *GoogleHandler {
	return &GoogleHandler{
		googleSvc: googleSvc,
	}
}

func (h *GoogleHandler) Autocomplete(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		response.HandleError(w, response.ErrBadRequest)
		return
	}

	predictions, err := h.googleSvc.Autocomplete(r.Context(), query)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.HandleSuccess(w, http.StatusOK, predictions)
}
