package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"log"
)

// Создание нового заказа
func CreateOrder(order models.Order) error {
	// Здесь можно добавить проверки и логику обработки заказа
	if err := dal.SaveOrder(order); err != nil {
		return err
	}
	log.Printf("Order created: %s", order.ID)
	return nil
}

// Дополнительные функции для обработки заказов (например, обновление статуса)

func GetAllOrders() ([]models.Order, error) {
	orders, err := dal.LoadOrders()
	if err != nil {
		log.Printf("Order created:")
		return nil, nil
	}
	return orders, nil
}
