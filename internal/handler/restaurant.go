package handler

import (
	"net/http"

	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/internal/service"
	"github.com/memsbdm/restaurant-api/internal/validation"
	"github.com/memsbdm/restaurant-api/pkg/keys"
)

type RestaurantHandler struct {
	restaurantSvc service.RestaurantService
}

func NewRestaurantHandler(restaurantSvc service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantSvc: restaurantSvc,
	}
}

type createRestaurantRequest struct {
	PlaceID string `json:"place_id" validate:"notblank"`
}

func (h *RestaurantHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request createRestaurantRequest

	if errs, err := validation.ValidateRequest(w, r, &request); err != nil || len(errs) != 0 {
		response.HandleValidationError(w, errs, err)
		return
	}

	userID, err := keys.GetUserIDFromContext(ctx)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	restaurant, err := h.restaurantSvc.Create(ctx, request.PlaceID, userID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.HandleSuccess(w, http.StatusCreated, restaurant)
}
