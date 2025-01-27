package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"log"
)

type MenuServiceInterface interface {
	AddMenuItem(menuItem models.MenuItem) error
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
