package handler

import (
	"net/http"
	"strings"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/dto"
	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/internal/service"
	"github.com/memsbdm/restaurant-api/internal/validation"
	"github.com/memsbdm/restaurant-api/pkg/contextkeys"
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

	createdUser, signedOAT, err := h.authSvc.Register(r.Context(), &dto.CreateUserDto{
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
			"access_token": signedOAT,
		})
		return
	}

	setAuthCookie(w, signedOAT, h.cfg.Env)
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

	fetchedUser, signedOAT, err := h.authSvc.Login(r.Context(), strings.TrimSpace(request.Email), request.Password)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	if IsMobileRequest(r) {
		response.HandleSuccess(w, http.StatusCreated, map[string]any{
			"user":         fetchedUser,
			"access_token": signedOAT,
		})
		return
	}

	setAuthCookie(w, signedOAT, h.cfg.Env)
	response.HandleSuccess(w, http.StatusCreated, fetchedUser)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	signedOAT, err := contextkeys.GetValueFromContext(r.Context(), contextkeys.SignedOATContextKey)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	err = h.authSvc.Logout(r.Context(), signedOAT)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.HandleSuccess(w, http.StatusNoContent, nil)
}

func setAuthCookie(w http.ResponseWriter, signedOAT, appEnv string) {
	cookie := &http.Cookie{
		Name:     "go-session",
		Value:    signedOAT,
		Path:     "/",
		HttpOnly: true,
		Secure:   appEnv == config.EnvProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 7,
	}

	http.SetCookie(w, cookie)
}
