package handlers

import (
	"encoding/json"
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
	LoginWithFirebase(w http.ResponseWriter, r *http.Request)
	RegistryUser(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	LogoutAllDevices(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	logger          ports.Logger
	authJWT         ports.AuthJWT
	userService     services.UserService
	authService     services.AuthTokenService
	firebaseService services.FirebaseService
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

	deviceInfo := r.Header.Get("X-Device-Info")
	if deviceInfo == "" {
		deviceInfo = "Unknown Device"
	}

	response, err := a.authService.Login(
		ctx,
		req.Email,
		req.Password,
		util.GetClientIP(r),
		r.Header.Get(middleware.UserAgentHeaderKey),
		deviceInfo,
	)

	if err != nil {
		a.logger.Error("Login failed", err)
		http_errors.ErrorHandler(w, err)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (a authHandler) LoginWithFirebase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.FirebaseLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("error decoding Firebase login request body: %v", err, err.Error())
		formatErr := apierrors.NewBadRequestRestError("invalid request body.")
		http_errors.ErrorHandler(w, formatErr)
		return
	}

	if err := req.Validate(); err != nil {
		http_errors.ErrorHandler(w, err)
		return
	}

	response, err := a.firebaseService.LoginWithFirebase(
		ctx,
		req.IDToken,
		util.GetClientIP(r),
		r.Header.Get(middleware.UserAgentHeaderKey),
		req.DeviceInfo,
	)

	if err != nil {
		a.logger.Error("Firebase login failed", err)
		http_errors.ErrorHandler(w, err)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
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
		DateBirth: *newUser.BirthDate,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)

}

func (a authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("error decoding request body: %v", err, err.Error())
		formatErr := apierrors.NewBadRequestRestError("invalid request body")
		http_errors.ErrorHandler(w, formatErr)
		return
	}

	accessToken, refreshToken, expiresIn, err := a.authService.RefreshAccessToken(ctx, req.RefreshToken)
	if err != nil {
		http_errors.ErrorHandler(w, err)
		return
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    string(core.TokenTypeBearer),
		ExpiresIn:    expiresIn,
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (a authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("error decoding request body for logout", err)
		formatErr := apierrors.NewBadRequestRestError("invalid request body")
		http_errors.ErrorHandler(w, formatErr)
		return
	}

	if err := a.authService.Logout(ctx, req.RefreshToken); err != nil {
		http_errors.ErrorHandler(w, err)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "logged out successfully",
	})
}

func (a authHandler) LogoutAllDevices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionKey, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("invalid session"))
		return
	}

	if err := a.authService.LogoutAllDevices(ctx, sessionKey.UserID); err != nil {
		http_errors.ErrorHandler(w, err)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "logged out from all devices successfully",
	})
}

func NewAuthHandler(
	logger ports.Logger,
	authJWT ports.AuthJWT,
	userService services.UserService,
	authService services.AuthTokenService,
	firebaseService services.FirebaseService,
) AuthHandler {
	return &authHandler{
		logger:          logger,
		authJWT:         authJWT,
		userService:     userService,
		authService:     authService,
		firebaseService: firebaseService,
	}
}
