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

func (m AuthMiddleware) AuthorizationTemporary(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	accessToken := r.Header.Get("x_access_token")
	if accessToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "unauthorized",
		})
		return
	}

	parts := strings.Split(accessToken, " ")
	if len(parts) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "unauthorized",
		})
		return
	}

	if parts[0] != "Bearer" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "unauthorized",
		})
		return
	}

	payload, err := m.encrypter.DecryptTemporary(parts[1])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "unauthorized",
		})
		return
	}

	req := r.WithContext(context.WithValue(r.Context(), UserIDKey{}, payload["sub"]))
	next(w, req)
}
