package dal

import (
	"encoding/json"
	"fmt"
	"hot-coffee/models"
	"os"
)

type OrderRepositoryInterface interface {
	CreateOrder(order models.Order) error
	LoadOrders() ([]models.Order, error)
}

type OrderRepositoryJSON struct {
	filePath string
}

func NewOrderRepositoryJSON(filePath string) *OrderRepositoryJSON {
	return &OrderRepositoryJSON{filePath: filePath}
}

const ordersFile = "data/orders.json"

// Сохранение заказа в JSON файл
func (r *OrderRepositoryJSON) CreateOrder(order models.Order) error {
	orders, err := r.LoadOrders()
	if err != nil {
		return err
	}

	orders = append(orders, order)

	return saveOrders(orders)
}

// Запись заказов в файл
func saveOrders(orders []models.Order) error {
	file, err := os.Create("data/orders.json")
	if err != nil {
		return fmt.Errorf("could not create orders file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Форматированный вывод JSON
	if err := encoder.Encode(orders); err != nil {
		return fmt.Errorf("could not encode orders to file: %v", err)
	}

	return nil
}

// Загрузка всех заказов из файла
func (r *OrderRepositoryJSON) LoadOrders() ([]models.Order, error) {
	file, err := os.Open(ordersFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var orders []models.Order
	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		return nil, err
	}

	return orders, nil
}
