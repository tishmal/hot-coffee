package service

import (
	"fmt"
	"hot-coffee/helper"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
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
	repository dal.OrderRepositoryInterface
}

func NewOrderService(repository dal.OrderRepositoryInterface) *OrderService {
	return &OrderService{repository: repository}
}

// Создание нового заказа
func (s *OrderService) CreateOrder(order models.Order) (models.Order, error) {
	newID := helper.GenerateID()

	for {
		if result, err := s.repository.GetOrderByID(strconv.Itoa(int(newID))); result == nil && err != nil {
			break
		}
		newID = helper.GenerateID()
	}

	order.ID = strconv.Itoa(int(newID))

	order.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	order.Status = "open"

	// Здесь можно добавить проверки и логику обработки заказа
	if err := s.repository.CreateOrder(order); err != nil {
		return order, err
	}
	log.Printf("Order created: %s", order.ID)
	return order, nil
}

// Дополнительные функции для обработки заказов (например, обновление статуса)

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
			return orders[i], nil
		}
	}
	return changeOrder, fmt.Errorf("Order with ID %s not found", id)
}

func (s *OrderService) CloseOrder(id string) (models.Order, error) {
	orders, err := s.repository.LoadOrders()
	if err != nil {
		return models.Order{}, fmt.Errorf("error reading all oreders %s: %v", id, err)
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
