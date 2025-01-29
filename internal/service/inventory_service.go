package service

import (
	"errors"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"log"
)

type InventoryServiceInterface interface {
	CreateInventory(Inventory models.InventoryItem) (models.InventoryItem, error)
	GetAllInventory() ([]models.InventoryItem, error)
}

type InventoryService struct {
	repository dal.InventoryRepositoryInterface
}

func NewInventoryService(_repository dal.InventoryRepositoryInterface) InventoryService {
	return InventoryService{repository: _repository}
}

func (s *InventoryService) CreateInventory(inventory models.InventoryItem) (models.InventoryItem, error) {
	if inventory.IngredientID == "" || inventory.Name == "" || inventory.Quantity == 0 || inventory.Unit == "" {
		return models.InventoryItem{}, errors.New("Invalid request body.")
	}

	if err := s.repository.CreateInventory(inventory); err != nil {
		return models.InventoryItem{}, errors.New("Inventory exists")
	}

	log.Printf("Inventory created: %s", inventory.IngredientID)
	return inventory, nil
}

func (h *InventoryService) GetAllInventory() ([]models.InventoryItem, error) {
	inventrories, err := h.repository.GetAllInventory()
	if err != nil {
		return nil, err
	}
	return inventrories, nil
}
