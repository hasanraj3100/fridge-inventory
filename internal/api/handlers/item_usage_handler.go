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
		switch err {
		case service.ErrItemNotFound:
			response.ResponseWithError(w, http.StatusNotFound, err.Error())
		case service.ErrUnauthorized:
			response.ResponseWithError(w, http.StatusForbidden, err.Error())
		case service.ErrInsufficientQuantity:
			response.ResponseWithError(w, http.StatusBadRequest, err.Error())
		default:
			log.Printf("CreateItemUsage error: %v", err)
			response.ResponseWithError(w, http.StatusInternalServerError, "Failed to create item usage")
		}
	}
}

func (iuh *ItemUsageHandler) GetItemUsageByUserID(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	result, err := iuh.itemUsageService.GetItemUsageByUserID(r.Context(), user.ID, page, pageSize)
	if err != nil {
		log.Printf("GetItemUsageByUserID error: %v", err)
		response.ResponseWithError(w, http.StatusInternalServerError, "Failed to get item usage")
		return
	}

	response.ResponseWithJSON(w, http.StatusOK, result)
}

func (iuh *ItemUsageHandler) UpdateItemUsage(w http.ResponseWriter, r *http.Request) {
	usageIDStr := r.PathValue("id")
	usageID, err := strconv.ParseInt(usageIDStr, 10, 64)
	if err != nil {
		response.ResponseWithError(w, http.StatusBadRequest, "Invalid usage ID")
		return
	}

	var req dto.ItemUsageUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&req)
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

	err = iuh.itemUsageService.UpdateItemUsage(r.Context(), user.ID, usageID, req)
	if err != nil {
		switch err {
		case service.ErrItemNotFound:
			response.ResponseWithError(w, http.StatusNotFound, err.Error())
		case service.ErrUnauthorized:
			response.ResponseWithError(w, http.StatusForbidden, err.Error())
		case service.ErrInsufficientQuantity:
			response.ResponseWithError(w, http.StatusBadRequest, err.Error())
		default:
			log.Printf("UpdateItemUsage error: %v", err)
			response.ResponseWithError(w, http.StatusInternalServerError, "Failed to update item usage")
		}
		return
	}

	response.ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "Item usage updated successfully"})
}

func (iuh *ItemUsageHandler) DeleteItemUsage(w http.ResponseWriter, r *http.Request) {
	usageIDStr := r.PathValue("id")
	usageID, err := strconv.ParseInt(usageIDStr, 10, 64)
	if err != nil {
		response.ResponseWithError(w, http.StatusBadRequest, "Invalid usage ID")
		return
	}

	user, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = iuh.itemUsageService.DeleteItemUsage(r.Context(), user.ID, usageID)
	if err != nil {
		switch err {
		case service.ErrItemNotFound:
			response.ResponseWithError(w, http.StatusNotFound, err.Error())
		case service.ErrUnauthorized:
			response.ResponseWithError(w, http.StatusForbidden, err.Error())
		default:
			log.Printf("DeleteItemUsage error: %v", err)
			response.ResponseWithError(w, http.StatusInternalServerError, "Failed to delete item usage")
		}
		return
	}

	response.ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "Item usage deleted successfully"})
}
