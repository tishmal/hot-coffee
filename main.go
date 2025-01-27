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

	// orderRepo := dal.NewOrderRepositoryJSON("")
	// orderService := service.NewOrderService(orderRepo)

	// orderHandler := handler.NewOrderHandler(*orderService)

	menuRepo := dal.NewMenuRepositoryJSON("")
	menuService := service.NewMenuService(menuRepo)

	menuHandler := handler.NewMenuHandler(menuService)

	// http.HandleFunc("/orders", handleOrders(orderHandler))
	// http.HandleFunc("/orders/", handleOrders(orderHandler))

	// http.HandleFunc("/orders/", orderHandler)
	http.HandleFunc("/menu", handleOrders(menuHandler))

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

// func handleOrders(orderHandler *handler.OrderHandler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		path := strings.Trim(r.URL.Path, "/")
// 		parts := strings.SplitN(path, "/", 2)

// 		switch r.Method {
// 		case http.MethodGet:
// 			if len(parts) == 1 {
// 				orderHandler.HandleGetAllOrders(w, r)
// 			} else if len(parts) == 2 {
// 				orderHandler.HandleGetOrderById(w, r, parts[1])
// 			} else {
// 				http.Error(w, "Not Found", http.StatusNotFound)
// 			}
// 		case http.MethodPost:
// 			if len(parts) == 1 {
// 				orderHandler.HandleCreateOrder(w, r)
// 			}
// 		// case http.MethodPut:
// 		// 	if len(parts) == 1 {
// 		// 		s.CreateBucket(w, r, parts[0])
// 		// 	} else if len(parts) == 2 {
// 		// 		s.PutObject(w, r, parts[0], parts[1])
// 		// 	} else {
// 		// 		utils.WriteErrorXML(w, "Bad Request", http.StatusBadRequest)
// 		// 	}
// 		case http.MethodDelete:
// 			if len(parts) == 2 {
// 				orderHandler.HandleDeleteOrder(w, r, parts[1])
// 			} else {
// 				http.Error(w, "Not Found", http.StatusNotFound)
// 			}
// 		default:
// 			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		}
// 	}
// }

func handleOrders(menuHandler *handler.MenuHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.SplitN(path, "/", 2)

		switch r.Method {
		case http.MethodPost:
			if len(parts) == 1 {
				menuHandler.HandleAddMenuItem(w, r)
			}
		// case http.MethodPut:
		// 	if len(parts) == 1 {
		// 		s.CreateBucket(w, r, parts[0])
		// 	} else if len(parts) == 2 {
		// 		s.PutObject(w, r, parts[0], parts[1])
		// 	} else {
		// 		utils.WriteErrorXML(w, "Bad Request", http.StatusBadRequest)
		// 	}

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
