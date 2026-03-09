package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

func (fih *FridgeItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	itemIDStr := r.PathValue("id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		response.ResponseWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req dto.FridgeItemUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := response.Validator.Struct(req); err != nil {
		niceErrors := response.FormatValidationError(err)
		response.ResponseWithValidationErrors(w, http.StatusBadRequest, "Validation failed", niceErrors)
		return
	}

	updatedItem, err := fih.fridgeItemService.UpdateItem(r.Context(), itemID, user.ID, req)
	if err != nil {
		if err.Error() == "item not found or does not belong to user" {
			response.ResponseWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		response.ResponseWithError(w, http.StatusInternalServerError, "Failed to update item")
		log.Printf("UpdateItem error: %v", err)
		return
	}

	response.ResponseWithJSON(w, http.StatusOK, updatedItem)
}

func (fih *FridgeItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	itemIDStr := r.PathValue("id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		response.ResponseWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	err = fih.fridgeItemService.DeleteItem(r.Context(), itemID, user.ID)
	if err != nil {
		if err.Error() == "item not found or does not belong to user" {
			response.ResponseWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		response.ResponseWithError(w, http.StatusInternalServerError, "Failed to delete item")
		log.Printf("DeleteItem error: %v", err)
		return
	}

	response.ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "Item deleted successfully"})
}
