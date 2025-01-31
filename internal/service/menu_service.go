package service

import (
	"errors"
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"hot-coffee/utils"
	"log"
)

type MenuServiceInterface interface {
	AddMenuItem(menuItem models.MenuItem) error
	GetAllMenuItems() ([]models.MenuItem, error)
	GetMenuItemByID(id string) (*models.MenuItem, error)
	DeleteMenuItemByID(id string) error
	UpdateMenu(id string, changeMenu models.MenuItem) (models.MenuItem, error)
}

type MenuService struct {
	repository dal.MenuRepositoryInterface
}

func NewMenuService(repository dal.MenuRepositoryInterface) *MenuService {
	return &MenuService{repository: repository}
}

func (m *MenuService) AddMenuItem(menuItem models.MenuItem) error {
	if err := utils.ValidateID(menuItem.ID); err != nil {
		return fmt.Errorf("invalid product ID: %v", err)
	}

	if err := utils.ValidateMenuItem(menuItem); err != nil {
		return err
	}

	if err := m.repository.AddMenuItem(menuItem); err != nil {
		return err
	}
	log.Printf("menu item added: %s", menuItem.ID)
	return nil
}

func (m *MenuService) GetAllMenuItems() ([]models.MenuItem, error) {
	items, err := m.repository.LoadMenuItems()
	if err != nil {
		log.Printf("could not load menu items: %v", err)
		return nil, fmt.Errorf("could not load menu items: %v", err)
	}
	return items, nil
}

func (m *MenuService) GetMenuItemByID(id string) (models.MenuItem, error) {
	if err := utils.ValidateID(id); err != nil {
		return models.MenuItem{}, fmt.Errorf("invalid menu ID: %v", err)
	}

	menuItems, err := m.repository.LoadMenuItems()
	if err != nil {
		return models.MenuItem{}, fmt.Errorf("could not load menu items: %v", err)
	}

	for _, item := range menuItems {
		if item.ID == id {
			return item, nil
		}
	}

	return models.MenuItem{}, fmt.Errorf("menu item with ID %s not found", id)
}

func (m *MenuService) DeleteMenuItemByID(id string) error {
	if err := utils.ValidateID(id); err != nil {
		return fmt.Errorf("invalid menu ID: %v", err)
	}

	menuItems, err := m.repository.LoadMenuItems()
	if err != nil {
		return fmt.Errorf("could not load menu items")
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

	if err := m.repository.SaveMenuItems(menuItems); err != nil {
		return fmt.Errorf("could not save menu items")
	}

	return nil
}

func (m *MenuService) UpdateMenu(id string, changeMenu models.MenuItem) (models.MenuItem, error) {
	if err := utils.ValidateID(id); err != nil {
		return models.MenuItem{}, err
	}

	if err := utils.ValidateMenuItem(changeMenu); err != nil {
		return models.MenuItem{}, err
	}

	menu, err := m.repository.LoadMenuItems()
	if err != nil {
		return models.MenuItem{}, errors.New("unable to load menu items")
	}

	for i := 0; i < len(menu); i++ {
		if menu[i].ID == id {
			if changeMenu.ID != menu[i].ID {
				return models.MenuItem{}, errors.New("cannot change the ID of the menu item")
			}
			menu[i].Name = changeMenu.Name
			menu[i].Ingredients = changeMenu.Ingredients
			menu[i].Description = changeMenu.Description
			menu[i].Price = changeMenu.Price
			m.repository.SaveMenuItems(menu)
			if err != nil {
				return models.MenuItem{}, fmt.Errorf("failed to save updated menu items: %v", err)
			}
			return menu[i], nil
		}
	}
	return models.MenuItem{}, errors.New("menu item not found")
}
