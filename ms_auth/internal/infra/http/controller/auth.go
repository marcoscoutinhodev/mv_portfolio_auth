package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/marcoscoutinhodev/ms_auth/internal/domain/user"
	"github.com/marcoscoutinhodev/ms_auth/internal/infra/http/mw"
	"github.com/marcoscoutinhodev/ms_auth/pkg"
)

type Auth struct {
	usecase user.UseCaseInterface
}

func NewAuth(usecase user.UseCaseInterface) *Auth {
	return &Auth{
		usecase: usecase,
	}
}

func (a Auth) SignUp(w http.ResponseWriter, r *http.Request) {
	var input user.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err,
		})
		return
	}

	input.Name = strings.ToUpper(strings.Join(strings.Fields(input.Name), " "))
	input.Email = strings.ToLower(strings.Join(strings.Fields(input.Email), " "))

	var errors []string

	if len(input.Name) < 5 {
		errors = append(errors, "name must contain at least 5 characters")
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		errors = append(errors, "invalid email format")
	}

	if !pkg.PasswordValidator(input.Password) {
		errors = append(errors, "password too weak")
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": errors,
		})
		return
	}

	output, err := a.usecase.Register(r.Context(), &input)
	if err != nil {
		fmt.Println("internal server error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}

func (a Auth) SignIn(w http.ResponseWriter, r *http.Request) {
	var input user.AuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err,
		})
		return
	}

	input.Email = strings.ToLower(strings.Join(strings.Fields(input.Email), " "))

	var errors []string

	if _, err := mail.ParseAddress(input.Email); err != nil {
		errors = append(errors, "poorly formatted email")
	}

	if !pkg.PasswordValidator(input.Password) {
		errors = append(errors, "poorly formatted password")
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": errors,
		})
		return
	}

	output, err := a.usecase.Auth(r.Context(), &input)
	if err != nil {
		fmt.Println("internal server error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}

func (a Auth) ForgottenPassword(w http.ResponseWriter, r *http.Request) {
	var input user.ForgottenPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err,
		})
		return
	}

	input.Email = strings.ToLower(strings.Join(strings.Fields(input.Email), " "))

	var errors []string

	if _, err := mail.ParseAddress(input.Email); err != nil {
		errors = append(errors, "poorly formatted email")
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": errors,
		})
		return
	}

	output, err := a.usecase.ForgottenPassword(r.Context(), &input)
	if err != nil {
		fmt.Println("internal server error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}

func (a Auth) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var input user.UpdatePasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err,
		})
		return
	}

	var errors []string

	if !pkg.PasswordValidator(input.Password) {
		errors = append(errors, "poorly formatted password")
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": errors,
		})
		return
	}

	input.UserID = r.Context().Value(mw.UserIDKey{}).(string)

	output, err := a.usecase.UpdatePassword(r.Context(), &input)
	if err != nil {
		fmt.Println("internal server error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}

func (a Auth) EmailConfirmationRequest(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")
	if _, err := mail.ParseAddress(email); err != nil {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "invalid email format",
			})
			return
		}
	}

	email = strings.ToLower(email)

	output, err := a.usecase.EmailConfirmationRequest(r.Context(), email)
	if err != nil {
		fmt.Println("internal server error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}

func (a Auth) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey{}).(string)

	output, err := a.usecase.ConfirmEmail(r.Context(), userID)
	if err != nil {
		fmt.Println("internal server error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}

func (a Auth) NewAccessToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(mw.UserIDKey{}).(string)

	output, err := a.usecase.NewAccessToken(r.Context(), userID)
	if err != nil {
		fmt.Println("internal server error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}
