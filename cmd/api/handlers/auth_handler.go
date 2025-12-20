// Package handlers implements HTTP handlers
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hasanraj3100/fridge-inventory/cmd/api/dto"
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
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := ah.userService.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			responseWithError(w, http.StatusConflict, err.Error())
		} else {
			fmt.Println("Register error:", err)
			responseWithError(w, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	token, err := ah.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
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

	token, err := ah.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			responseWithError(w, http.StatusUnauthorized, err.Error())
		} else {
			responseWithError(w, http.StatusInternalServerError, "Failed to login")
			fmt.Println("Login error:", err)
		}
		return
	}

	if err := validate.Struct(req); err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := ah.userService.GetUserByEmail(r.Context(), req.Email)
	if err != nil || user == nil {
		responseWithError(w, http.StatusInternalServerError, "Failed to retrieve user details")
		fmt.Println("GetUserByEmail error:", err)
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
