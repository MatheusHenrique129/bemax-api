package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MatheusHenrique129/bemax-api/internal/adapters/constants"
	"github.com/MatheusHenrique129/bemax-api/internal/adapters/handlers/middleware"
	"github.com/MatheusHenrique129/bemax-api/internal/core"
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
	"github.com/MatheusHenrique129/bemax-api/internal/util"
	"github.com/MatheusHenrique129/bemax-api/pkg/http_errors"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	RegistryUser(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	logger      ports.Logger
	authJWT     ports.AuthJWT
	userService services.UserService
	authService services.AuthTokenService
}

func (a authHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("error decoding request body: %v", err, err.Error())
		formatErr := apierrors.NewBadRequestRestError("invalid request body.")
		http_errors.ErrorHandler(w, formatErr)
		return
	}

	accessToken, refreshToken, err := a.authService.Login(
		ctx,
		req.Email,
		req.Password,
		util.GetClientIP(r),
		r.Header.Get(middleware.UserAgentHeaderKey),
	)

	if err != nil {
		a.logger.Error("Login failed", err)
		http_errors.ErrorHandler(w, err)
		return
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    string(core.TokenTypeBearer),
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.Header().Set(middleware.AuthHeaderKey, fmt.Sprintf("Bearer %s", accessToken))

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (a authHandler) RegistryUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("error decoding request body: %v", err, err.Error())
		formatErr := apierrors.NewBadRequestRestError("error decoding request body.")
		http_errors.ErrorHandler(w, formatErr)
		return
	}

	// Validate input data via DTO
	if _, err := req.Validate(); err != nil {
		a.logger.Error("Request validation failed %v", err, err.Error())
		http_errors.ErrorHandler(w, err)
		return
	}

	newUser, err := a.userService.CreateUser(ctx, req)
	if err != nil {
		http_errors.ErrorHandler(w, err)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)

	response := dto.UserRegisterResponse{
		Email:     newUser.Email,
		FullName:  newUser.FullName,
		CPF:       newUser.CPF,
		Phone:     newUser.Phone,
		DateBirth: newUser.BirthDate,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)

}

func NewAuthHandler(
	logger ports.Logger,
	authJWT ports.AuthJWT,
	userService services.UserService,
	authService services.AuthTokenService,
) AuthHandler {
	return &authHandler{
		logger:      logger,
		authJWT:     authJWT,
		userService: userService,
		authService: authService,
	}
}
