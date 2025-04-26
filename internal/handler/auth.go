package handler

import (
	"net/http"
	"strings"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/dto"
	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/internal/service"
	"github.com/memsbdm/restaurant-api/internal/validation"
	"github.com/memsbdm/restaurant-api/pkg/keys"
)

type AuthHandler struct {
	cfg     *config.App
	authSvc service.AuthService
}

func NewAuthHandler(cfg *config.App, authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{
		cfg:     cfg,
		authSvc: authSvc,
	}
}

type registerUserRequest struct {
	Name     string `validate:"notblank,max=50"`
	Email    string `validate:"notblank,email"`
	Password string `validate:"notblank,min=8"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request registerUserRequest

	if errs, err := validation.ValidateRequest(w, r, &request); err != nil || len(errs) != 0 {
		response.HandleValidationError(w, errs, err)
		return
	}

	createdUser, oat, err := h.authSvc.Register(r.Context(), &dto.CreateUser{
		Name:     strings.TrimSpace(request.Name),
		Email:    strings.TrimSpace(request.Email),
		Password: request.Password,
	})
	if err != nil {
		response.HandleError(w, err)
		return
	}

	if IsMobileRequest(r) {
		response.HandleSuccess(w, http.StatusCreated, map[string]any{
			"user":         createdUser,
			"access_token": oat,
		})
		return
	}

	SetAuthCookie(w, oat, h.cfg.Env)
	response.HandleSuccess(w, http.StatusCreated, createdUser)
}

type loginUserRequest struct {
	Email    string `validate:"notblank,email"`
	Password string `validate:"notblank"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request loginUserRequest

	if errs, err := validation.ValidateRequest(w, r, &request); err != nil || len(errs) != 0 {
		response.HandleValidationError(w, errs, err)
		return
	}

	loginResponse, oat, err := h.authSvc.Login(r.Context(), strings.TrimSpace(request.Email), request.Password)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	if IsMobileRequest(r) {
		response.HandleSuccess(w, http.StatusCreated, map[string]any{
			"user":         loginResponse.User,
			"access_token": oat,
			"restaurants":  loginResponse.Restaurants,
		})
		return
	}

	SetAuthCookie(w, oat, h.cfg.Env)
	response.HandleSuccess(w, http.StatusCreated, loginResponse)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	oat, err := keys.GetValueFromContext(r.Context(), keys.AuthOATContextKey)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	err = h.authSvc.Logout(r.Context(), oat)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	if !IsMobileRequest(r) {
		clearAuthCookie(w, h.cfg.Env)
	}

	response.HandleSuccess(w, http.StatusNoContent, nil)
}
