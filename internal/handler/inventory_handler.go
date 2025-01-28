package handler

import "hot-coffee/internal/service"

type InventoryHandlerInterface interface {
}

type InventoryHandler struct {
	_inventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) InventoryHandler {
	return InventoryHandler{_inventoryService: inventoryService}
}
