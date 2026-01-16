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

type FridgeItemHandler struct {
	fridgeItemService service.FridgeItemService
}

func NewFridgeItemHandler(fridgeItemService service.FridgeItemService) *FridgeItemHandler {
	return &FridgeItemHandler{
		fridgeItemService: fridgeItemService,
	}
}

func (fih *FridgeItemHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	var req dto.FridgeItemAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	req.UserID = user.ID

	createdItem, err := fih.fridgeItemService.AddItem(r.Context(), req)
	if err != nil {
		response.ResponseWithError(w, http.StatusInternalServerError, "Failed to add fridge item")
		log.Printf("AddItem error: %v", err)
		return
	}
	response.ResponseWithJSON(w, http.StatusCreated, createdItem)
}

func (fih *FridgeItemHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	items, err := fih.fridgeItemService.GetByUserID(r.Context(), user.ID)
	if err != nil {
		response.ResponseWithError(w, http.StatusInternalServerError, "Failed to get items")
		log.Printf("GetItem error: %v", err)
		return
	}

	response.ResponseWithJSON(w, http.StatusOK, items)
}
