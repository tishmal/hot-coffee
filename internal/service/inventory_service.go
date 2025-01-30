package service

import (
	"errors"
	"log"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type InventoryServiceInterface interface {
	CreateInventory(Inventory models.InventoryItem) (models.InventoryItem, error)
	GetAllInventory() ([]models.InventoryItem, error)
	GetInventoryByID(id string) (models.InventoryItem, error)
}

type InventoryService struct {
	repository dal.InventoryRepositoryInterface
}

func NewInventoryService(_repository dal.InventoryRepositoryInterface) InventoryService {
	return InventoryService{repository: _repository}
}

func (s *InventoryService) CreateInventory(inventory models.InventoryItem) (models.InventoryItem, error) {
	if inventory.IngredientID == "" || inventory.Name == "" || inventory.Quantity == 0 || inventory.Unit == "" {
		return models.InventoryItem{}, errors.New("invalid request body")
	}

	if err := s.repository.CreateInventory(inventory); err != nil {
		return models.InventoryItem{}, errors.New("inventory exists")
	}

	log.Printf("Inventory created: %s", inventory.IngredientID)
	return inventory, nil
}

func (s *InventoryService) GetAllInventory() ([]models.InventoryItem, error) {
	inventrories, err := s.repository.GetAllInventory()
	if err != nil {
		return nil, err
	}
	return inventrories, nil
}

func (s *InventoryService) GetInventoryByID(id string) (models.InventoryItem, error) {
	inventory, err := s.repository.GetAllInventory()
	if err != nil {
		return models.InventoryItem{}, errors.New("get all inventory")
	}

	if len(inventory) > 0 {
		for i := 0; i < len(inventory); i++ {
			if inventory[i].IngredientID == id {
				return inventory[i], nil
			}
		}
	} else {
		return models.InventoryItem{}, errors.New("inventory is empty")
	}
	return models.InventoryItem{}, errors.New("invalid product ID in inventory items")
}
