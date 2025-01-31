package dal

import (
	"encoding/json"
	"fmt"
	"hot-coffee/models"
	"os"
	"path/filepath"
)

type InventoryRepositoryInterface interface {
	CreateInventory(inventory models.InventoryItem) error
	SaveInventory(inventories []models.InventoryItem) error
	GetAllInventory() ([]models.InventoryItem, error)
}

type InventoryRepositoryJSON struct {
	filePath string
}

func NewInventoryRepositoryJSON(filepath string) InventoryRepositoryJSON {
	return InventoryRepositoryJSON{filePath: filepath}
}

func (r InventoryRepositoryJSON) CreateInventory(inventory models.InventoryItem) error {
	inventories, err := r.GetAllInventory()
	if err != nil {
		return err
	}

	inventories = append(inventories, inventory)

	return r.SaveInventory(inventories)
}

func (r InventoryRepositoryJSON) SaveInventory(inventories []models.InventoryItem) error {
	filePath := filepath.Join(r.filePath, "inventory.json")
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not create orders file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Форматированный вывод JSON
	if err := encoder.Encode(inventories); err != nil {
		return fmt.Errorf("could not encode orders to file: %v", err)
	}

	return nil
}

func (r InventoryRepositoryJSON) GetAllInventory() ([]models.InventoryItem, error) {
	filePath := filepath.Join(r.filePath, "inventory.json")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var inventories []models.InventoryItem
	if err := json.NewDecoder(file).Decode(&inventories); err != nil {
		return nil, err
	}

	return inventories, nil
}
