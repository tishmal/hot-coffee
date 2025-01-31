package handler

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"hot-coffee/utils"
	"net/http"
)

type MenuHandlerInterface interface {
	HandleAddMenuItem(w http.ResponseWriter, r *http.Request)
	HandleGetAllMenuItems(w http.ResponseWriter, r *http.Request)
	HandleGetMenuItemById(w http.ResponseWriter, r *http.Request, menuID string)
	HandleDeleteMenuItemById(w http.ResponseWriter, r *http.Request, menuID string)
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
		utils.ErrorInJSON(w, http.StatusBadRequest, fmt.Errorf("invalid JSON format: %v", err))
		return
	}

	if err := utils.ValidateMenuItem(NewMenuItem); err != nil {
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

	utils.ResponseInJSON(w, 201, NewMenuItem)
}

func (m *MenuHandler) HandleGetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := m.menuService.GetAllMenuItems()
	if err != nil {
		utils.ErrorInJSON(w, http.StatusInternalServerError, err)
		return
	}

	if len(items) == 0 {
		utils.ResponseInJSON(w, 200, []models.MenuItem{})
		return
	}

	utils.ResponseInJSON(w, 200, items)
}

func (m *MenuHandler) HandleGetMenuItemById(w http.ResponseWriter, r *http.Request, menuID string) {
	if err := utils.ValidateID(menuID); err != nil {
		utils.ErrorInJSON(w, http.StatusBadRequest, fmt.Errorf("invalid menu ID: %v", err))
		return
	}

	item, err := m.menuService.GetMenuItemByID(menuID)
	if err != nil || &item == nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	}

	utils.ResponseInJSON(w, 200, item)
}

func (m *MenuHandler) HandleDeleteMenuItemById(w http.ResponseWriter, r *http.Request, menuID string) {
	if err := utils.ValidateID(menuID); err != nil {
		utils.ErrorInJSON(w, http.StatusBadRequest, fmt.Errorf("invalid menu ID: %v", err))
		return
	}

	err := m.menuService.DeleteMenuItemByID(menuID)
	if err != nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (m *MenuHandler) HandleUpdateMenu(w http.ResponseWriter, r *http.Request, menuID string) {
	if err := utils.ValidateID(menuID); err != nil {
		utils.ErrorInJSON(w, http.StatusBadRequest, err)
		return
	}

	var changeMenu models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&changeMenu); err != nil {
		utils.ErrorInJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.ValidateMenuItem(changeMenu); err != nil {
		utils.ErrorInJSON(w, http.StatusBadRequest, err)
		return
	}

	if menu, err := m.menuService.UpdateMenu(menuID, changeMenu); err != nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	} else {
		utils.ResponseInJSON(w, 200, menu)
	}
}
