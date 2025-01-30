package dal

import (
	"encoding/json"
	"fmt"
	"hot-coffee/models"
	"io/ioutil"
	"log"
	"os"
)

type OrderRepositoryInterface interface {
	CreateOrder(order models.Order) error
	LoadOrders() ([]models.Order, error)
	GetOrderByID(id string) (*models.Order, error)
	DeleteOrder(id string) (*models.Order, error)
	UpdateOrder(id string, changeOrder models.Order) ([]models.Order, error)
	SaveOrders(orders []models.Order) error
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

	return r.SaveOrders(orders)
}

// Запись заказов в файл
func (r *OrderRepositoryJSON) SaveOrders(orders []models.Order) error {
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

func (r *OrderRepositoryJSON) GetOrderByID(id string) (*models.Order, error) {
	// Открываем файл с данными
	data, err := ioutil.ReadFile("data/orders.json")
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}

	// Слайс для хранения всех заказов
	var orders []models.Order

	// Парсим JSON из файла в структуру
	err = json.Unmarshal(data, &orders)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Проверяем, что заказы есть в данных
	if len(orders) > 0 {
		for i := 0; i < len(orders); i++ {
			if orders[i].ID == id {
				return &orders[i], nil
			}
		}
	}
	return nil, fmt.Errorf("Order with ID %s not found", id)
}

func (r *OrderRepositoryJSON) DeleteOrder(id string) (*models.Order, error) {
	// Открываем файл с данными
	data, err := ioutil.ReadFile("data/orders.json")
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}

	// Слайс для хранения всех заказов
	var orders []models.Order

	// Парсим JSON из файла в структуру
	err = json.Unmarshal(data, &orders)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
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

func (r *OrderRepositoryJSON) UpdateOrder(id string, changeOrder models.Order) ([]models.Order, error) {
	orders, err := r.LoadOrders()
	if err != nil {
		return orders, err
	}
	return orders, nil
}
