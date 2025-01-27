package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"net/http"
)

type MenuHandlerInterface interface {
	HandleAddMenuItem(w http.ResponseWriter, r *http.Request)
	HandleGetAllMenuItems(w http.ResponseWriter, r *http.Request)
	HandleGetMenuItemById(w http.ResponseWriter, r *http.Request, orderID string)
	HandleDeleteMenuItemById(w http.ResponseWriter, r *http.Request, orderID string)
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

func (m *MenuHandler) HandleGetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := m.menuService.GetAllMenuItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(items)
}

func (m *MenuHandler) HandleGetMenuItemById(w http.ResponseWriter, r *http.Request, orderID string) {
	item, err := m.menuService.GetMenuItemByID(orderID)
	if err != nil || &item == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (m *MenuHandler) HandleDeleteMenuItemById(w http.ResponseWriter, r *http.Request, orderID string) {
	err := m.menuService.DeleteMenuItemByID(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
