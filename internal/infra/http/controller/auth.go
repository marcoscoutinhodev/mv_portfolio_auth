package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/marcoscoutinhodev/mv_chat/internal/domain/user"
	"github.com/marcoscoutinhodev/mv_chat/pkg"
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
		errors = append(errors, "invalid email")
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

	output := a.usecase.Register(r.Context(), &input)

	if output.StatusCode == http.StatusInternalServerError {
		fmt.Println("internal server error:", output.Error)
		w.WriteHeader(output.StatusCode)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}
