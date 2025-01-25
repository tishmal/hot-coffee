package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"log"
)

type OrderServiceInterface interface {
	CreateOrder(order models.Order) error
	GetAllOrders() ([]models.Order, error)
}

type OrderService struct {
	repository dal.OrderRepositoryInterface
}

func NewOrderService(repository dal.OrderRepositoryInterface) *OrderService {
	return &OrderService{repository: repository}
}

// Создание нового заказа
func (s *OrderService) CreateOrder(order models.Order) error {
	// Здесь можно добавить проверки и логику обработки заказа
	if err := s.repository.CreateOrder(order); err != nil {
		return err
	}
	log.Printf("Order created: %s", order.ID)
	return nil
}

// Дополнительные функции для обработки заказов (например, обновление статуса)

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	orders, err := s.repository.LoadOrders()
	if err != nil {
		log.Printf("Order created:")
		return nil, nil
	}
	return orders, nil
}
