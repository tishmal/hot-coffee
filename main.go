package main

import (
	"flag"
	"fmt"
	"hot-coffee/add"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	. "hot-coffee/models"
	"log"
	"net/http"
	"strings"
)

var (
	orders         = make(map[string]Order)
	menuItems      = make(map[string]MenuItem)
	inventoryItems = make(map[string]InventoryItem)
)

var orderHandler *handler.OrderHandler

func main() {
	port := flag.Int("port", 8080, "Port number to listen on")
	dir := flag.String("dir", "data", "Base directory for storage")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		add.PrintUsage()
		return
	}

	orderRepo := dal.NewOrderRepositoryJSON(*dir)
	orderService := service.NewOrderService(orderRepo)
	orderHandler = handler.NewOrderHandler(*orderService)

	http.HandleFunc("/orders", handleOrders(orderHandler))
	http.HandleFunc("/orders/", handleOrders(orderHandler))
	// http.HandleFunc("/orders", orderHandler.HandleGetAllOrders)
	// http.HandleFunc("/orders/", orderHandler)
	// http.HandleFunc("/menu", menuHandler)
	// http.HandleFunc("/menu/", menuHandler)
	// http.HandleFunc("/inventory", inventoryHandler)
	// http.HandleFunc("/inventory/", inventoryHandler)

	addr := fmt.Sprintf(":%d", *port)

	// Запуск браузера
	// go add.OpenBrowser(addr)

	// Запуск HTTP сервера
	log.Printf("The server is running on the port %s...\n", addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleOrders(orderHandler *handler.OrderHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := strings.TrimPrefix(r.URL.Path, "/orders/")

		switch {
		case orderID == "" && r.Method == http.MethodGet:
			// Handle Get All Orders
			orderHandler.HandleGetAllOrders(w, r)
		case orderID == "" && r.Method == http.MethodPost:
			// Handle Create New Order
			orderHandler.HandleCreateOrder(w, r)
		case orderID != "" && r.Method == http.MethodGet:
			// Handle Get Order By ID
			orderHandler.HandleGetOrdersByID(w, r, orderID)
		case orderID != "" && r.Method == http.MethodPut:
			// Handle Update Order By ID
			orderHandler.HandleUpdateOrder(w, r, orderID)
		case orderID != "" && r.Method == http.MethodDelete:
			// Handle Delete Order By ID
			orderHandler.HandleDeleteOrder(w, r, orderID)
		case orderID != "" && r.Method == http.MethodPost:
			// Handle Close Order By ID
			orderHandler.HandleCloseOrder(w, r, orderID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// func menuHandler(w http.ResponseWriter, r *http.Request) {
// 	menuID := strings.TrimPrefix(r.URL.Path, "/menu/")

// 	switch r.Method {
// 	case http.MethodPut:
// 		updateMenuItem(w, r, menuID)
// 	case http.MethodGet:
// 		if menuID == "" {
// 			getAllmenuItems(w)
// 		} else {
// 			getMenuItemByID(w, r, menuID)
// 		}
// 	case http.MethodDelete:
// 		deleteMenuItem(w, menuID)
// 	case http.MethodPost:
// 		addMenuItem(w, r)
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// func inventoryHandler(w http.ResponseWriter, r *http.Request) {
// 	inventoryID := strings.TrimPrefix(r.URL.Path, "/inventory/")

// 	switch r.Method {
// 	case http.MethodPut:
// 		updateInventoryItem(w, r, inventoryID)
// 	case http.MethodGet:
// 		if menuID == "" {
// 			getAllInventoryItems(w)
// 		} else {
// 			getInventoryItemByID(w, r, inventoryID)
// 		}
// 	case http.MethodDelete:
// 		deleteInventoryItem(w, inventoryID)
// 	case http.MethodPost:
// 		addInventoryItem(w, r)
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// // Обработчик запросов
// func handleRequests(w http.ResponseWriter, r *http.Request) {
// 	// Устанавливаем заголовки для правильного отображения HTML
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	// Чтение HTML файла
// 	html, err := ioutil.ReadFile("index.html")
// 	if err != nil {
// 		http.Error(w, "Failed to load HTML file", http.StatusInternalServerError)
// 		return
// 	}

// 	// Вставляем путь из URL в HTML-страницу
// 	username := r.URL.Path[1:]
// 	if username == "" {
// 		username = "Guest"
// 	}

// 	// Заменяем метку <username> на значение из URL
// 	htmlStr := strings.Replace(string(html), "<username>", username, 1)

// 	// Отправляем страницу клиенту
// 	fmt.Fprint(w, htmlStr)
// }
