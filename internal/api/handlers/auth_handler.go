// Package handlers implements HTTP handlers
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hasanraj3100/fridge-inventory/internal/api/dto"
	"github.com/hasanraj3100/fridge-inventory/internal/service"
)

var validate = validator.New()

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validate.Struct(req); err != nil {
		niceErrors := formatValidationError(err)
		responseWithValidationErrors(w, http.StatusBadRequest, "Validation failed", niceErrors)
		return
	}

	err := ah.userService.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			responseWithError(w, http.StatusConflict, err.Error())
		} else {
			log.Printf("Register error: %v", err)
			responseWithError(w, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	token, user, err := ah.userService.Login(r.Context(), dto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Token generation error: %v", err)
		responseWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := dto.AuthResponse{
		Token: token,
		User: dto.UserDetail{
			ID:       user.ID,
			Username: user.UserName,
			Email:    user.Email,
		},
	}

	responseWithJSON(w, http.StatusCreated, response)
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validate.Struct(req); err != nil {
		niceErrors := formatValidationError(err)
		responseWithValidationErrors(w, http.StatusBadRequest, "Validation failed", niceErrors)
		return
	}

	token, user, err := ah.userService.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			responseWithError(w, http.StatusUnauthorized, err.Error())
		} else {
			log.Printf("Login error: %v", err)
			responseWithError(w, http.StatusInternalServerError, "Failed to login")
		}
		return
	}

	response := dto.AuthResponse{
		Token: token,
		User: dto.UserDetail{
			ID:       user.ID,
			Username: user.UserName,
			Email:    user.Email,
		},
	}

	responseWithJSON(w, http.StatusOK, response)
}

func responseWithJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		fmt.Println("Error encoding response:", err)
	}
}

func responseWithError(w http.ResponseWriter, statusCode int, message string) {
	responseWithJSON(w, statusCode, map[string]string{"error": message})
}

func responseWithValidationErrors(w http.ResponseWriter, statusCode int, message string, details any) {
	responseWithJSON(w, statusCode, map[string]any{
		"error":   message,
		"details": details,
	})
}

func formatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if vdErrors, ok := err.(validator.ValidationErrors); ok {
		for _, f := range vdErrors {
			switch f.Tag() {
			case "required":
				errors[f.Field()] = "this field is required"
			case "email":
				errors[f.Field()] = "invalid email format"
			case "min":
				errors[f.Field()] = fmt.Sprintf("must be at least %s characters", f.Param())
			default:
				errors[f.Field()] = fmt.Sprintf("failed validation on %s", f.Tag())
			}
		}
	}
	return errors
}
