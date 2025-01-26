package main

import (
	"flag"
	"fmt"
	"hot-coffee/add"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"log"
	"net/http"
	"strings"
)

func main() {
	port := flag.Int("port", 8080, "Port number to listen on")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		add.PrintUsage()
		return
	}

	orderRepo := dal.NewOrderRepositoryJSON("")
	orderService := service.NewOrderService(orderRepo)

	orderHandler := handler.NewOrderHandler(*orderService)

	http.HandleFunc("/orders", handleOrders(orderHandler))
	http.HandleFunc("/orders/", handleOrders(orderHandler))

	// http.HandleFunc("/orders/", orderHandler)
	// http.HandleFunc("/menu", menuHandler)
	// http.HandleFunc("/menu/", menuHandler)
	// http.HandleFunc("/inventory", inventoryHandler)
	// http.HandleFunc("/inventory/", inventoryHandler)

	addr := fmt.Sprintf(":%d", *port)

	// Запуск браузера
	go add.OpenBrowser(addr)

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
			orderHandler.HandleGetOrderById(w, r, orderID)
		case orderID != "" && r.Method == http.MethodPut:
			// Handle Update Order By ID
			// orderHandler.HandleUpdateOrder(w, r, orderID)
		case orderID != "" && r.Method == http.MethodDelete:
			// Handle Delete Order By ID
			// orderHandler.HandleDeleteOrder(w, r, orderID)
		case orderID != "" && r.Method == http.MethodPost:
			// Handle Close Order By ID
			// orderHandler.HandleCloseOrder(w, r, orderID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

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
