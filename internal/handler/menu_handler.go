package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"net/http"
)

type MenuHandlerInterface interface {
	HandleAddMenuItem(w http.ResponseWriter, r *http.Request)
}

type MenuHandler struct {
	menuService *service.MenuService
}

func NewMenuHandler(menuService *service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

func (m *MenuHandler) HandleAddMenuItem(w http.ResponseWriter, r *http.Request) {
	var NewMenuItem models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&NewMenuItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	menuItem := models.MenuItem{
		ID:          NewMenuItem.ID,
		Name:        NewMenuItem.Name,
		Description: NewMenuItem.Description,
		Price:       NewMenuItem.Price,
		Ingredients: NewMenuItem.Ingredients,
	}

	if err := m.menuService.AddMenuItem(menuItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(NewMenuItem)
}
