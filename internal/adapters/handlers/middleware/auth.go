package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	auth "github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services"
	"github.com/MatheusHenrique129/bemax-api/pkg/http_errors"
)

var (
	AuthHeaderKey      = "Authorization"
	UserAgentHeaderKey = "User-Agent"
)

type contextKey string

const sessionContextKey contextKey = "sessionKey"

type SessionKey struct {
	Claims *auth.Claims
}

type AuthMiddleware interface {
	AuthenticateRequest(next http.Handler) http.Handler
}
type authMiddleware struct {
	jwtPort      ports.AuthJWT
	logger       ports.Logger
	tokenService services.AuthTokenService
}

func (a authMiddleware) AuthenticateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := validateAuthHeader(r.Header.Get(AuthHeaderKey))
		if err != nil {
			http_errors.ErrorHandler(w, err)
			return
		}

		session, err := a.authenticateUserSession(token)
		if err != nil {
			a.logger.Error(err.Error(), err)
			http_errors.ErrorHandler(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), sessionContextKey, SessionKey{Claims: session})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a authMiddleware) authenticateUserSession(tokenString string) (*auth.Claims, apierrors.RestError) {
	claims, err := a.tokenService.ValidateAccessToken(tokenString)
	if err != nil {
		a.logger.Error("Error validating token", err)
		return nil, err
	}

	a.logger.Debug(fmt.Sprintf("User %s authenticated, claims: %v", claims.UserID.String(), claims))
	return claims, nil
}

func validateAuthHeader(authHeader string) (string, apierrors.RestError) {
	if authHeader == "" {
		return "", apierrors.NewUnauthorizedRestError(fmt.Sprintf("token is required - missing '%s' in headers", AuthHeaderKey))
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", apierrors.NewUnauthorizedRestError("authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

func GetClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	val, ok := ctx.Value(sessionContextKey).(SessionKey)
	return val.Claims, ok && val.Claims != nil
}

func NewAuthMiddleware(
	logger ports.Logger,
	jwtPort ports.AuthJWT,
	tokenService services.AuthTokenService,
) AuthMiddleware {
	return authMiddleware{
		logger:       logger,
		jwtPort:      jwtPort,
		tokenService: tokenService,
	}
}
