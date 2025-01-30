package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"hot-coffee/helper"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
)

func main() {
	port := flag.Int("port", 8080, "Port number to listen on")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		helper.PrintUsage()
		return
	}
	// 1
	orderRepo := dal.NewOrderRepositoryJSON("")
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(*orderService)
	// 2
	inventoryRepo := dal.NewInventoryRepositoryJSON("")
	inventoryService := service.NewInventoryService(inventoryRepo)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	http.HandleFunc("/orders", handleRequestsOrders(orderHandler))
	http.HandleFunc("/orders/", handleRequestsOrders(orderHandler))
	// http.HandleFunc("/menu", menuHandler)
	// http.HandleFunc("/menu/", menuHandler)
	http.HandleFunc("/inventory", handleRequestsInventory(inventoryHandler))
	http.HandleFunc("/inventory/", handleRequestsInventory(inventoryHandler))

	addr := fmt.Sprintf(":%d", *port)

	// Запуск браузера
	go helper.OpenBrowser(addr)

	// Запуск HTTP сервера
	log.Printf("The server is running on the port %s...\n", addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequestsInventory(inventoryHandler handler.InventoryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.SplitN(path, "/", 3)

		switch r.Method {
		case http.MethodGet:
			if len(parts) == 1 {
				inventoryHandler.HandleGetAllInventory(w, r)
			} else if len(parts) == 2 {
				inventoryHandler.HandleGetInventoryById(w, r, parts[1])
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		case http.MethodPost:
			if len(parts) == 1 {
				inventoryHandler.HandleCreateInventory(w, r)
			} else {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
		case http.MethodPut:
			if len(parts) == 2 {
				// orderHandler.HandleUpdateOrder(w, r, parts[1])
			} else {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
		case http.MethodDelete:
			if len(parts) == 2 {
				// orderHandler.HandleDeleteOrder(w, r, parts[1])
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleRequestsOrders(orderHandler *handler.OrderHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.SplitN(path, "/", 3)

		switch r.Method {
		case http.MethodGet:
			if len(parts) == 1 {
				orderHandler.HandleGetAllOrders(w, r)
			} else if len(parts) == 2 {
				orderHandler.HandleGetOrderById(w, r, parts[1])
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		case http.MethodPost:
			if len(parts) == 1 {
				orderHandler.HandleCreateOrder(w, r)
			} else if len(parts) == 3 && parts[2] == "close" {
				orderHandler.HandleCloseOrder(w, r, parts[1])
			} else {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
		case http.MethodPut:
			if len(parts) == 2 {
				orderHandler.HandleUpdateOrder(w, r, parts[1])
			} else {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
		case http.MethodDelete:
			if len(parts) == 2 {
				orderHandler.HandleDeleteOrder(w, r, parts[1])
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
