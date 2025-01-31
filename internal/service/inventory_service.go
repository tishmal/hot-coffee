package service

import (
	"errors"
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"hot-coffee/utils"
	"log"
)

type InventoryServiceInterface interface {
	CreateInventory(Inventory models.InventoryItem) (models.InventoryItem, error)
	GetAllInventory() ([]models.InventoryItem, error)
	GetInventoryByID(id string) (models.InventoryItem, error)
	DeleteInventoryItemByID(is string) error
	UpdateInventoryItem(inventoryItemID string, changedInventoryItem models.InventoryItem) (models.InventoryItem, error)
}

type InventoryService struct {
	repository dal.InventoryRepositoryInterface
}

func NewInventoryService(_repository dal.InventoryRepositoryInterface) InventoryService {
	return InventoryService{repository: _repository}
}

func (s *InventoryService) CreateInventory(inventory models.InventoryItem) (models.InventoryItem, error) {
	if err := utils.IsValidName(inventory.Name); err != nil {
		return models.InventoryItem{}, err
	}

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

func (h *InventoryService) DeleteInventoryItemByID(id string) error {
	inventoryItems, err := h.repository.GetAllInventory()
	if err != nil {
		return err
	}

	indexToDelete := -1
	for i, item := range inventoryItems {
		if item.IngredientID == id {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		return fmt.Errorf("inventory item with ID %s not found", id)
	}

	inventoryItems = append(inventoryItems[:indexToDelete], inventoryItems[indexToDelete+1:]...)

	return h.repository.SaveInventory(inventoryItems)
}

func (h *InventoryService) UpdateInventoryItem(inventoryItemID string, changedInventoryItem models.InventoryItem) (models.InventoryItem, error) {
	if changedInventoryItem.Name == "" || changedInventoryItem.Quantity == 0 || changedInventoryItem.Unit == "" {
		return models.InventoryItem{}, errors.New("invalid request body")
	}

	inventoryItem, err := h.repository.GetAllInventory()
	if err != nil {
		return models.InventoryItem{}, errors.New("invalid load inventory items")
	}

	for i := 0; i < len(inventoryItem); i++ {
		if inventoryItem[i].IngredientID == inventoryItemID {
			inventoryItem[i].Name = changedInventoryItem.Name
			inventoryItem[i].Quantity = changedInventoryItem.Quantity
			inventoryItem[i].Unit = changedInventoryItem.Unit
			h.repository.SaveInventory(inventoryItem)
			if changedInventoryItem.IngredientID != inventoryItem[i].IngredientID {
				return models.InventoryItem{}, errors.New("cannot change ID")
			} else {
				return inventoryItem[i], nil
			}
		}
	}
	return models.InventoryItem{}, errors.New("invalid ID in inventory items")
}
