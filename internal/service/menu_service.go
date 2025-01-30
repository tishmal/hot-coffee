package service

import (
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"log"
)

type MenuServiceInterface interface {
	AddMenuItem(menuItem models.MenuItem) error
	GetAllMenuItems() ([]models.MenuItem, error)
	GetMenuItemByID(id string) (*models.MenuItem, error)
	DeleteMenuItemByID(id string) error
}

type MenuService struct {
	repository dal.MenuRepositoryInterface
}

func NewMenuService(repository dal.MenuRepositoryInterface) *MenuService {
	return &MenuService{repository: repository}
}

func (m *MenuService) AddMenuItem(menuItem models.MenuItem) error {
	if err := m.repository.AddMenuItem(menuItem); err != nil {
		return err
	}
	log.Printf("menu item added: %s", menuItem.ID)
	return nil
}

func (m *MenuService) GetAllMenuItems() ([]models.MenuItem, error) {
	items, err := m.repository.LoadMenuItems()
	if err != nil {
		log.Printf("items list created:")
		return nil, fmt.Errorf("could not load menu items: %v", err)
	}
	return items, nil
}

func (m *MenuService) GetMenuItemByID(id string) (models.MenuItem, error) {
	menuItems, err := m.repository.LoadMenuItems()
	if err != nil {
		return models.MenuItem{}, err
	}

	for _, item := range menuItems {
		if item.ID == id {
			return item, nil
		}
	}

	return models.MenuItem{}, fmt.Errorf("menu item with ID %s not found", id)
}

func (m *MenuService) DeleteMenuItemByID(id string) error {
	menuItems, err := m.repository.LoadMenuItems()
	if err != nil {
		return err
	}

	indexToDelete := -1
	for i, item := range menuItems {
		if item.ID == id {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		return fmt.Errorf("menu item with ID %s not found", id)
	}

	menuItems = append(menuItems[:indexToDelete], menuItems[indexToDelete+1:]...)

	return m.repository.SaveMenuItems(menuItems)
}
