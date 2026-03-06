package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hasanraj3100/fridge-inventory/internal/api/dto"
	"github.com/hasanraj3100/fridge-inventory/internal/api/response"
	"github.com/hasanraj3100/fridge-inventory/internal/middleware"
	"github.com/hasanraj3100/fridge-inventory/internal/service"
)

type ItemUsageHandler struct {
	itemUsageService service.ItemUsageService
}

func NewItemUsageHandler(itemUsageService service.ItemUsageService) *ItemUsageHandler {
	return &ItemUsageHandler{
		itemUsageService: itemUsageService,
	}
}

func (iuh *ItemUsageHandler) CreateItemUsage(w http.ResponseWriter, r *http.Request) {
	var req dto.ItemUsageCreateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := response.Validator.Struct(req); err != nil {
		niceErrors := response.FormatValidationError(err)
		response.ResponseWithValidationErrors(w, http.StatusBadRequest, "Validation failed", niceErrors)
		return
	}
	user, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = iuh.itemUsageService.CreateItemUsage(r.Context(), user.ID, req)
	if err != nil {
		log.Printf("CreateItemUsage error: %v", err)
		response.ResponseWithError(w, http.StatusInternalServerError, "Failed to create item usage")
		return
	}
	response.ResponseWithJSON(w, http.StatusCreated, map[string]string{"message": "Item usage recorded successfully"})
}
