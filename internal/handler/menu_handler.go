package handler

import (
	"encoding/json"
	"net/http"

	"hot-coffee/internal/service"
	"hot-coffee/models"
	"hot-coffee/utils"
)

type MenuHandlerInterface interface {
	HandleAddMenuItem(w http.ResponseWriter, r *http.Request)
	HandleGetAllMenuItems(w http.ResponseWriter, r *http.Request)
	HandleGetMenuItemById(w http.ResponseWriter, r *http.Request, orderID string)
	HandleDeleteMenuItemById(w http.ResponseWriter, r *http.Request, orderID string)
	HandleUpdateMenu(w http.ResponseWriter, r *http.Request, menuID string)
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
		utils.ErrorInJSON(w, http.StatusBadRequest, err)
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
		utils.ErrorInJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseInJSON(w, NewMenuItem)
}

func (m *MenuHandler) HandleGetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := m.menuService.GetAllMenuItems()
	if err != nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
	}

	utils.ResponseInJSON(w, items)
}

func (m *MenuHandler) HandleGetMenuItemById(w http.ResponseWriter, r *http.Request, orderID string) {
	item, err := m.menuService.GetMenuItemByID(orderID)
	if err != nil || &item == nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	}

	utils.ResponseInJSON(w, item)
}

func (m *MenuHandler) HandleDeleteMenuItemById(w http.ResponseWriter, r *http.Request, orderID string) {
	err := m.menuService.DeleteMenuItemByID(orderID)
	if err != nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (m *MenuHandler) HandleUpdateMenu(w http.ResponseWriter, r *http.Request, menuID string) {
	var changeMenu models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&changeMenu); err != nil {
		utils.ErrorInJSON(w, http.StatusBadRequest, err)
		return
	}

	if menu, err := m.menuService.UpdateMenu(menuID, changeMenu); err != nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	} else {
		utils.ResponseInJSON(w, menu)
	}
}
