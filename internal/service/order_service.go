package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"hot-coffee/helper"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"hot-coffee/utils"
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

	err = utils.ValidateOrder(menu, order)
	if err != nil {
		return models.Order{}, err
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
	order, err := s.GetOrderByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete order with ID %s: %v", id, err)
	}
	
	return order, nil
}
return nil, fmt.Errorf("failed to delete order with ID %s: %v", id, err)
func (r *OrderRepositoryJSON) DeleteOrder(id string) (*models.Order, error) {
	orders, err := r.LoadOrders()
	if err != nil {
		return &models.Order{}, err
	}

	for i := 0; i < len(orders); i++ {
		if orders[i].ID == id {
			deletedOrder := orders[i]

			orders = append(orders[:i], orders[i+1:]...)

			updatedData, err := json.MarshalIndent(orders, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("Error marshaling updated orders: %v", err)
			}

			err = ioutil.WriteFile("data/orders.json", updatedData, os.ModePerm)
			if err != nil {
				return nil, fmt.Errorf("Error writing updated file: %v", err)
			}

			return &deletedOrder, nil
		}
	}
	return nil, fmt.Errorf("Order with ID %s not found", id)
}

func (s *OrderService) UpdateOrder(id string, changeOrder models.Order) (models.Order, error) {
	if changeOrder.CustomerName == "" || changeOrder.Items == nil {
		return models.Order{}, errors.New("invalid request body")
	}

	orders, err := s.repository.LoadOrders()
	if err != nil {
		return changeOrder, fmt.Errorf("error reading all oreders %s: %v", id, err)
	}

	menu, err := s.menuService.repository.LoadMenuItems()
	if err != nil {
		return models.Order{}, err
	}

	err = utils.ValidateOrder(menu, changeOrder)
	if err != nil {
		return models.Order{}, err
	}

	for i := 0; i < len(orders); i++ {
		if orders[i].ID == changeOrder.ID && changeOrder.Status == "closed" {
			return models.Order{}, fmt.Errorf("order is closed")
		}
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

	orderId, err := s.repository.GetOrderByID(id)
	if err != nil {
		return models.Order{}, fmt.Errorf("Failed to retrieve order by ID%v", id)
	}

	if orderId.Status == "closed" {
		return models.Order{}, fmt.Errorf("opration not allowed")
	}

	menu, err := s.menuService.repository.LoadMenuItems()
	if err != nil {
		return models.Order{}, err
	}

	menuMap := make(map[string]models.MenuItem)
	for _, items := range menu {
		menuMap[items.ID] = items
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

	for _, items := range orderId.Items {
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

	for i := 0; i < len(orders); i++ {
		if orders[i].ID == id {
			orders[i].Status = "closed"
			s.repository.SaveOrders(orders)
			return orders[i], nil
		}
	}
	return models.Order{}, fmt.Errorf("Order with ID %s not found", id)
}
