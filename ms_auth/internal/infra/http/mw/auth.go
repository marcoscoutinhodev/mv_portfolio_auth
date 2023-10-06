package mw

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/marcoscoutinhodev/ms_auth/internal/infra/adapter"
)

type UserIDKey struct{}

type AuthMiddleware struct {
	encrypter adapter.EncrypterInterface
}

func NewAuthMiddleware(encrypter adapter.EncrypterInterface) *AuthMiddleware {
	return &AuthMiddleware{
		encrypter: encrypter,
	}
}

func (m AuthMiddleware) preValidation(token string) (bool, string) {
	if token == "" {
		return false, ""
	}

	parts := strings.Split(token, " ")
	if len(parts) != 2 {
		return false, ""
	}

	if parts[0] != "Bearer" {
		return false, ""
	}

	return true, parts[1]
}

func (m AuthMiddleware) Authorization(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	accessToken := r.Header.Get("x_access_token")

	if isValid, token := m.preValidation(accessToken); isValid {
		payload, err := m.encrypter.Decrypt(token)
		if err == nil {
			req := r.WithContext(context.WithValue(r.Context(), UserIDKey{}, payload["sub"]))
			next(w, req)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": "unauthorized",
	})
}

func (m AuthMiddleware) AuthorizationTemporary(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	accessToken := r.Header.Get("x_access_token")

	if isValid, token := m.preValidation(accessToken); isValid {
		payload, err := m.encrypter.DecryptTemporary(token)
		if err == nil {
			req := r.WithContext(context.WithValue(r.Context(), UserIDKey{}, payload["sub"]))
			next(w, req)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": "unauthorized",
	})
}
