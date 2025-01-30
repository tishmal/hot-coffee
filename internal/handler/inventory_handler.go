package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"hot-coffee/utils"
	"net/http"
)

type InventoryHandlerInterface interface {
	HandleCreateInventory(w http.ResponseWriter, r *http.Request)
	HandleGetAllInventory(w http.ResponseWriter, r *http.Request)
	HandleGetInventoryById(w http.ResponseWriter, r *http.Request, id string)
	HandleDeleteInventoryItem(w http.ResponseWriter, r *http.Request, inventoryItemID string)
	HandleUpdateInventoryItem(w http.ResponseWriter, r *http.Request, inventoryItemID string)
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
		utils.ErrorInJSON(w, http.StatusBadRequest, err)
		return
	}

	// Вызываем сервис для создания заказа
	if inventory, err := h.inventoryService.CreateInventory(newInventory); err != nil {
		utils.ErrorInJSON(w, http.StatusBadRequest, err)
		return
	} else {
		utils.ResponseInJSON(w, inventory)
	}
}

func (h *InventoryHandler) HandleGetAllInventory(w http.ResponseWriter, r *http.Request) {
	inventories, err := h.inventoryService.GetAllInventory()
	if err != nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
	}

	utils.ResponseInJSON(w, inventories)
}

func (h *InventoryHandler) HandleGetInventoryById(w http.ResponseWriter, r *http.Request, id string) {
	inventory, err := h.inventoryService.GetInventoryByID(id)
	if err != nil {
		utils.ErrorInJSON(w, 404, err)
		return
	}

	utils.ResponseInJSON(w, inventory)
}

func (h *InventoryHandler) HandleDeleteInventoryItem(w http.ResponseWriter, r *http.Request, inventoryItemID string) {
	err := h.inventoryService.DeleteInventoryItemByID(inventoryItemID)
	if err != nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h *InventoryHandler) HandleUpdateInventoryItem(w http.ResponseWriter, r *http.Request, inventoryItemID string) {
	var changedInventoryItem models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&changedInventoryItem); err != nil {
		utils.ErrorInJSON(w, http.StatusBadRequest, err)
		return
	}

	if item, err := h.inventoryService.UpdateInventoryItem(inventoryItemID, changedInventoryItem); err != nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	} else {
		utils.ResponseInJSON(w, item)
	}
}
