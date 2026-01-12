// Package handlers implements HTTP handlers
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/hasanraj3100/fridge-inventory/internal/api/dto"
	"github.com/hasanraj3100/fridge-inventory/internal/api/response"
	"github.com/hasanraj3100/fridge-inventory/internal/service"
)

var validate = response.Validator

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validate.Struct(req); err != nil {
		niceErrors := response.FormatValidationError(err)
		response.ResponseWithValidationErrors(w, http.StatusBadRequest, "Validation failed", niceErrors)
		return
	}

	err := ah.userService.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			response.ResponseWithError(w, http.StatusConflict, err.Error())
		} else {
			log.Printf("Register error: %v", err)
			response.ResponseWithError(w, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	token, user, err := ah.userService.Login(r.Context(), dto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Token generation error: %v", err)
		response.ResponseWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	res := dto.AuthResponse{
		Token: token,
		User: dto.UserDetail{
			ID:       user.ID,
			Username: user.UserName,
			Email:    user.Email,
		},
	}

	response.ResponseWithJSON(w, http.StatusCreated, res)
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validate.Struct(req); err != nil {
		niceErrors := response.FormatValidationError(err)
		response.ResponseWithValidationErrors(w, http.StatusBadRequest, "Validation failed", niceErrors)
		return
	}

	token, user, err := ah.userService.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.ResponseWithError(w, http.StatusUnauthorized, err.Error())
		} else {
			log.Printf("Login error: %v", err)
			response.ResponseWithError(w, http.StatusInternalServerError, "Failed to login")
		}
		return
	}

	res := dto.AuthResponse{
		Token: token,
		User: dto.UserDetail{
			ID:       user.ID,
			Username: user.UserName,
			Email:    user.Email,
		},
	}

	response.ResponseWithJSON(w, http.StatusOK, res)
}
