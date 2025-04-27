package handler

import (
	"net/http"

	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/internal/service"
	"github.com/memsbdm/restaurant-api/internal/validation"
	"github.com/memsbdm/restaurant-api/pkg/keys"
)

type MenuHandler struct {
	menuSvc service.MenuService
}

func NewMenuHandler(menuSvc service.MenuService) *MenuHandler {
	return &MenuHandler{
		menuSvc: menuSvc,
	}
}

type CreateMenuRequest struct {
	Name string `json:"name"`
}

func (h *MenuHandler) Create(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := keys.GetRestaurantIDFromContext(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	var request CreateMenuRequest
	if errs, err := validation.ValidateRequest(w, r, &request); err != nil || len(errs) != 0 {
		response.HandleValidationError(w, errs, err)
		return
	}

	menu, err := h.menuSvc.Create(r.Context(), request.Name, restaurantID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.HandleSuccess(w, http.StatusCreated, menu)
}
