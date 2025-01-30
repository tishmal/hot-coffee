package handler

import (
	"encoding/json"
	"net/http"

	"hot-coffee/internal/service"
	"hot-coffee/models"
	"hot-coffee/utils"
)

type InventoryHandlerInterface interface {
	HandleCreateInventory(w http.ResponseWriter, r *http.Request)
	HandleGetAllInventory(w http.ResponseWriter, r *http.Request)
}

type InventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(_inventoryService service.InventoryService) InventoryHandler {
	return InventoryHandler{inventoryService: _inventoryService}
}

func (h *InventoryHandler) HandleCreateInventory(w http.ResponseWriter, r *http.Request) {
	var newInventory models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&newInventory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вызываем сервис для создания заказа
	if inventory, err := h.inventoryService.CreateInventory(newInventory); err != nil {
		utils.ErrorInJSON(w, http.StatusBadRequest, err)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(inventory)
	}
}

func (h *InventoryHandler) HandleGetAllInventory(w http.ResponseWriter, r *http.Request) {
	inventories, err := h.inventoryService.GetAllInventory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(inventories)
}

func (h *InventoryHandler) HandleGetInventoryById(w http.ResponseWriter, r *http.Request, id string) {
	inventory, err := h.inventoryService.GetInventoryByID(id)
	if err != nil {
		utils.ErrorInJSON(w, 404, err)
		return
	}

	utils.ResponseInJSON(w, inventory)
}
