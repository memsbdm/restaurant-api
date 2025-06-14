package handler

import (
	"net/http"

	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/internal/service"
	"github.com/memsbdm/restaurant-api/pkg/keys"
)

type VerifyEmailHandler struct {
	userService service.UserService
}

func NewVerifyEmailHandler(userSvc service.UserService) *VerifyEmailHandler {
	return &VerifyEmailHandler{
		userService: userSvc,
	}
}

func (h *VerifyEmailHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	spt := r.URL.Query().Get("token")
	if spt == "" {
		response.HandleError(w, response.ErrBadRequest)
		return
	}

	updatedUser, err := h.userService.VerifyEmail(r.Context(), spt)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.HandleSuccess(w, http.StatusOK, updatedUser)
}

func (h *VerifyEmailHandler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	userID, err := keys.GetUserIDFromContext(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	err = h.userService.ResendVerificationEmail(r.Context(), userID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.HandleSuccess(w, http.StatusOK, nil)
}
