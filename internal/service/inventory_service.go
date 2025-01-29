package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"log"
)

type InventoryServiceInterface interface {
	CreateInventory(Inventory models.InventoryItem) (models.InventoryItem, error)
}

type InventoryService struct {
	repository dal.InventoryRepositoryInterface
}

func NewInventoryService(_repository dal.InventoryRepositoryInterface) InventoryService {
	return InventoryService{repository: _repository}
}

func (s *InventoryService) CreateInventory(inventory models.InventoryItem) (models.InventoryItem, error) {
	if err := s.repository.CreateInventory(inventory); err != nil {
		return inventory, err
	}
	log.Printf("Inventory created: %s", inventory.IngredientID)
	return inventory, nil
}
