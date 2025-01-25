package main

import (
	. "hot-coffee/models"
	"log"
	"net/http"
	"strings"
)

var orders = make(map[string]Order)

func main() {
	// Регистрация маршрутов для создания заказов
	http.HandleFunc("/orders", orderHandler)
	http.HandleFunc("/orders/", orderHandler)
	http.HandleFunc("/menu", menuHandler)
	http.HandleFunc("/menu/", menuHandler)
	http.HandleFunc("/inventory", inventoryHandler)
	http.HandleFunc("/inventory/", inventoryHandler)

	// Запуск HTTP-сервера
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimPrefix(r.URL.Path, "/orders/")

	switch r.Method {
	case http.MethodPut:
		updateOrder(w, r, orderID)
	case http.MethodGet:
		if orderID == "" {
			getAllOrders(w)
		} else {
			getOrderByID(w, r, orderID)
		}
	case http.MethodDelete:
		deleteOrder(w, orderID)
	case http.MethodPost:
		if orderID == "" {
			createOrder{w, r}
		} else {
			closeOrder(w, r, orderID)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	menuID := strings.TrimPrefix(r.URL.Path, "/menu/")

	switch r.Method {
	case http.MethodPut:
		updateMenuItem(w, r, menuID)
	case http.MethodGet:
		if menuID == "" {
			getAllmenuItems(w)
		} else {
			getMenuItemByID(w, r, menuID)
		}
	case http.MethodDelete:
		deleteMenuItem(w, menuID)
	case http.MethodPost:
		addMenuItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func inventoryHandler(w http.ResponseWriter, r *http.Request) {
	inventoryID := strings.TrimPrefix(r.URL.Path, "/inventory/")

	switch r.Method {
	case http.MethodPut:
		updateInventoryItem(w, r, inventoryID)
	case http.MethodGet:
		if menuID == "" {
			getAllInventoryItems(w)
		} else {
			getInventoryItemByID(w, r, inventoryID)
		}
	case http.MethodDelete:
		deleteInventoryItem(w, inventoryID)
	case http.MethodPost:
		addInventoryItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
