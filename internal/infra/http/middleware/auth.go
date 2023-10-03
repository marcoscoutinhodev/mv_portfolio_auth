package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type UserIDKey struct{}

type Middleware struct {
	encrypter Encrypter
}

func NewMiddleware(encrypter Encrypter) *Middleware {
	return &Middleware{
		encrypter: encrypter,
	}
}

func (m Middleware) Authorization(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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

	payload, err := m.encrypter.Decrypt(parts[1])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "unauthorized",
		})
		return
	}

	p := payload["payload"].(map[string]interface{})
	userID := p["sub"].(string)

	req := r.WithContext(context.WithValue(r.Context(), UserIDKey{}, userID))
	next(w, req)
}
