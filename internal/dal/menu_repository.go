package dal

import (
	"encoding/json"
	"fmt"
	"hot-coffee/models"
	"io/ioutil"
	"os"
)

type MenuRepositoryInterface interface {
	AddMenuItem(menuItem models.MenuItem) error
	LoadMenuItems() ([]models.MenuItem, error)
	SaveMenuItems(menuItems []models.MenuItem) error
}

type MenuRepositoryJSON struct {
	filePath string
}

func NewMenuRepositoryJSON(filePath string) *MenuRepositoryJSON {
	return &MenuRepositoryJSON{filePath: filePath}
}

func (r *MenuRepositoryJSON) AddMenuItem(menuItem models.MenuItem) error {
	menuItems, err := r.LoadMenuItems()
	if err != nil {
		return err
	}

	menuItems = append(menuItems, menuItem)

	return saveMenuItems(menuItems)
}

func saveMenuItems(menuItems []models.MenuItem) error {
	file, err := os.OpenFile("data/menu_items.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not create menu file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Форматированный вывод JSON
	if err := encoder.Encode(menuItems); err != nil {
		return fmt.Errorf("could not encode menu items to file: %v", err)
	}

	return nil
}

func (m *MenuRepositoryJSON) LoadMenuItems() ([]models.MenuItem, error) {
	file, err := os.Open("data/menu_items.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []models.MenuItem{}, nil
		}
		return nil, fmt.Errorf("could not open menu file: %v", err)
	}
	defer file.Close()

	var menuItems []models.MenuItem
	if err := json.NewDecoder(file).Decode(&menuItems); err != nil {
		return nil, fmt.Errorf("could not decode menu: %v", err)
	}

	return menuItems, nil
}

func (r *MenuRepositoryJSON) SaveMenuItems(menuItems []models.MenuItem) error {
	// Запись обновленных данных в файл
	updatedData, err := json.MarshalIndent(menuItems, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling menu items: %v", err)
	}

	// Сохранение в файл
	err = ioutil.WriteFile("data/menu_items.json", updatedData, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing to menu items file: %v", err)
	}

	return nil
}
