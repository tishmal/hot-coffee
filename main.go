package main

import (
	"flag"
	"fmt"
	"hot-coffee/helper"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"hot-coffee/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.Int("port", 8080, "Port number to listen on")
	help := flag.Bool("help", false, "Show help")
	dir := flag.String("dir", "data", "Directory path for storing data")
	flag.Parse()

	if *help {
		helper.PrintUsage()
		return
	}
	// 2
	inventoryRepo := dal.NewInventoryRepositoryJSON(*dir)
	inventoryService := service.NewInventoryService(inventoryRepo)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	menuRepo := dal.NewMenuRepositoryJSON(*dir)
	menuService := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)

	// 1
	orderRepo := dal.NewOrderRepositoryJSON(*dir)
	orderService := service.NewOrderService(orderRepo, *menuService, inventoryService)
	orderHandler := handler.NewOrderHandler(*orderService)

	reportService := service.NewReportService(*menuService, *orderService)
	reportHandler := handler.NewReportHandler(reportService)

	http.HandleFunc("/orders", handleRequestsOrders(orderHandler))
	http.HandleFunc("/orders/", handleRequestsOrders(orderHandler))

	http.HandleFunc("/menu", handleMenu(menuHandler))
	http.HandleFunc("/menu/", handleMenu(menuHandler))
	//
	http.HandleFunc("/inventory", handleRequestsInventory(inventoryHandler))
	http.HandleFunc("/inventory/", handleRequestsInventory(inventoryHandler))

	http.HandleFunc("/reports/", handleRequestsReports(reportHandler))

	addr := fmt.Sprintf(":%d", *port)

	if err := utils.IsValidName(*dir); err != nil {
		log.Fatal("Invalid directory name: %w", err)
	}

	if err := os.MkdirAll(*dir, 0o755); err != nil {
		log.Fatalf("Error creating data directory: %v\n", err)
		os.Exit(1)
	}

	if *port < 0 || *port > 65535 {
		log.Fatal("Invalid port number")
	}

	// Запуск браузера
	go helper.OpenBrowser(addr)

	// Запуск HTTP сервера
	log.Printf("Server running on port %s with BaseDir %s\n", addr, *dir)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequestsReports(reportHandler handler.ReportHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.SplitN(path, "/", 2)

		switch r.Method {
		case http.MethodGet:
			if len(parts) == 2 && parts[1] == "total-sales" {
				reportHandler.HandleGetTotalSales(w, r)
			} else if len(parts) == 2 && parts[1] == "popular-items" {
				reportHandler.HandleGetPopulatItem(w, r)
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
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
				inventoryHandler.HandleUpdateInventoryItem(w, r, parts[1])
			} else {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
		case http.MethodDelete:
			if len(parts) == 2 {
				inventoryHandler.HandleDeleteInventoryItem(w, r, parts[1])
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

func handleMenu(menuHandler *handler.MenuHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.SplitN(path, "/", 2)

		switch r.Method {
		case http.MethodPost:
			if len(parts) == 1 {
				menuHandler.HandleAddMenuItem(w, r)
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		case http.MethodGet:
			if len(parts) == 1 {
				menuHandler.HandleGetAllMenuItems(w, r)
			} else if len(parts) == 2 {
				menuHandler.HandleGetMenuItemById(w, r, parts[1])
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		case http.MethodPut:
			if len(parts) == 2 {
				menuHandler.HandleUpdateMenu(w, r, parts[1])
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		case http.MethodDelete:
			if len(parts) == 2 {
				menuHandler.HandleDeleteMenuItemById(w, r, parts[1])
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
