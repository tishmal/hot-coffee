package service

import (
	"errors"
	"fmt"
	"hot-coffee/helper"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"hot-coffee/utils"
	"log"
	"strconv"
	"time"
)

type OrderServiceInterface interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetAllOrders() ([]models.Order, error)
	GetOrderByID(id string) (*models.Order, error)
	DeleteOrder(id string) (*models.Order, error)
	UpdateOrder(id string) (models.Order, error)
	CloseOrder(orderID string) (models.Order, error)
}

type OrderService struct {
	repository       dal.OrderRepositoryInterface
	menuService      MenuService
	inventoryService InventoryService
}

func NewOrderService(_repository dal.OrderRepositoryInterface, _menuService MenuService, _inventoryService InventoryService) *OrderService {
	return &OrderService{
		repository:       _repository,
		menuService:      _menuService,
		inventoryService: _inventoryService,
	}
}

func (s *OrderService) CreateOrder(order models.Order) (models.Order, error) {
	if err := utils.IsValidName(order.CustomerName); err != nil {
		return models.Order{}, err
	}

	newID := helper.GenerateID()

	for {
		if result, err := s.repository.GetOrderByID("order" + strconv.Itoa(int(newID))); result == nil && err != nil {
			break
		}
		newID = helper.GenerateID()
	}

	menu, err := s.menuService.repository.LoadMenuItems()
	if err != nil {
		return models.Order{}, err
	}

	menuMap := make(map[string]models.MenuItem)
	for _, items := range menu {
		menuMap[items.ID] = items
	}

	// Прежде чем перейдём в наличие ингридиентов, проверяем меню:
	// Валидация на соответсвие отсылаемого запросом и исполняемого заказа с тем что есть в меню
	// 1. Добавляем в массивы данные которые совпали с теми, что хранятся в menu.json
	var idshki []string
	for i := 0; i < len(menu); i++ {
		for _, item := range order.Items {
			if item.ProductID == menu[i].ID {
				idshki = append(idshki, item.ProductID)
			}
		}
	}
	// 2. Перебираем заказ и пробиваем на валидацию
	for _, item := range order.Items {
		err := utils.ValidateQuantity(float64(item.Quantity)) // преждевременная валидация на большие и отрицательные цифрры она не плоха
		if err != nil {
			return models.Order{}, err
		}

		if len(idshki) == 0 {
			return models.Order{}, fmt.Errorf("Invalid product ID: %s. Product not found in the menu.", item.ProductID)
		}

		for i := 0; i < len(idshki); i++ {
			if item.ProductID != idshki[i] {
				return models.Order{}, fmt.Errorf("Invalid product ID: %s. Product not found in the menu.", item.ProductID)
			}
		}
	}

	inventory, err := s.inventoryService.GetAllInventory()
	if err != nil {
		return models.Order{}, fmt.Errorf("Failed to retrieve inventory")
	}

	var newDataMenu []models.MenuItem

	ingredientMap := make(map[string]models.InventoryItem)
	for _, items := range inventory {
		ingredientMap[items.IngredientID] = items
	}

	for _, items := range order.Items {
		for i := 0; i < items.Quantity; i++ {
			if item, exists := menuMap[items.ProductID]; exists {
				newDataMenu = append(newDataMenu, item)
			}
		}
		if err := utils.ValidateQuantity(float64(items.Quantity)); err != nil {
			return models.Order{}, err
		}
	}

	for _, items := range newDataMenu {
		for _, ingredient := range items.Ingredients {
			if item, exist := ingredientMap[ingredient.IngredientID]; exist {
				fmt.Printf("Checking ingredient ID: %v, required: %v, available: %v\n",
					ingredient.IngredientID, ingredient.Quantity, item.Quantity)
				if ingredient.Quantity > item.Quantity {
					return models.Order{}, fmt.Errorf("not enough quantity for ingredient ID %v: required %v, available %v", ingredient.IngredientID, ingredient.Quantity, item.Quantity)
				}
				item.Quantity -= ingredient.Quantity
				ingredientMap[ingredient.IngredientID] = item
			}
		}
	}

	for ingredientID, item := range ingredientMap {
		if _, err := s.inventoryService.UpdateInventoryItem(ingredientID, item); err != nil {
			return models.Order{}, fmt.Errorf("Failed to update inventory for ingredientID %v", ingredientID)
		}
	}

	order.ID = "order" + strconv.Itoa(int(newID))
	order.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	order.Status = "open"

	if err := s.repository.CreateOrder(order); err != nil {
		return order, err
	}
	log.Printf("Order created: %s", order.ID)
	return order, nil
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	orders, err := s.repository.LoadOrders()
	if err != nil {
		log.Printf("error get all orders!")
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {
	order, err := s.repository.GetOrderByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}
	return order, nil
}

func (s *OrderService) DeleteOrder(id string) (*models.Order, error) {
	order, err := s.repository.DeleteOrder(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete order with ID %s: %v", id, err)
	}
	return order, nil
}

func (s *OrderService) UpdateOrder(id string, changeOrder models.Order) (models.Order, error) {
	if changeOrder.CustomerName == "" || changeOrder.Items == nil || changeOrder.Status == "" {
		return models.Order{}, errors.New("invalid request body")
	}

	orders, err := s.repository.UpdateOrder(id, changeOrder)
	if err != nil {
		return changeOrder, fmt.Errorf("error reading all oreders %s: %v", id, err)
	}

	for i := 0; i < len(orders); i++ {
		if orders[i].ID == id {
			orders[i].CustomerName = changeOrder.CustomerName
			orders[i].CreatedAt = time.Now().UTC().Format(time.RFC3339)
			orders[i].Items = changeOrder.Items
			s.repository.SaveOrders(orders)
			if changeOrder.ID != orders[i].ID {
				return models.Order{}, errors.New("cannot change ID")
			} else {
				return orders[i], nil
			}
		}
	}
	return changeOrder, fmt.Errorf("Order with ID %s not found", id)
}

func (s *OrderService) CloseOrder(id string) (models.Order, error) {
	orders, err := s.repository.LoadOrders()
	if err != nil {
		return models.Order{}, fmt.Errorf("Order with ID %s not found", id)
	}

	for i := 0; i < len(orders); i++ {
		if orders[i].ID == id {
			orders[i].Status = "closed"
			s.repository.SaveOrders(orders)
			return orders[i], nil
		}
	}
	return models.Order{}, fmt.Errorf("Order with ID %s not found", id)
}
