package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"net/http"
)

type InventoryHandlerInterface interface {
	HandleCreateInventory(w http.ResponseWriter, r *http.Request)
}

type InventoryHandler struct {
	_inventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) InventoryHandler {
	return InventoryHandler{_inventoryService: inventoryService}
}

func (h *InventoryHandler) HandleCreateInventory(w http.ResponseWriter, r *http.Request) {
	var newInventory models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&newInventory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вызываем сервис для создания заказа
	if inventory, err := h._inventoryService.CreateInventory(newInventory); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(inventory)
	}
}
