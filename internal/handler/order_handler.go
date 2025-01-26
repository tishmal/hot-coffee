package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"net/http"
)

type OrderHandlerInterface interface {
	HandleCreateOrder(w http.ResponseWriter, r *http.Request)
	HandleGetAllOrders(w http.ResponseWriter, r *http.Request)
	HandleGetOrderById(w http.ResponseWriter, r *http.Request)
	HandleDeleteOrder(w http.ResponseWriter, r *http.Request)
}

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// Обработчик для создания нового заказа
func (h *OrderHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вызываем сервис для создания заказа
	if err := h.orderService.CreateOrder(newOrder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
}

// Дополнительные обработчики для получения, обновления и удаления заказов

func (h *OrderHandler) HandleGetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) HandleGetOrderById(w http.ResponseWriter, r *http.Request, orderID string) {
	order, err := h.orderService.GetOrderByID(orderID)
	if err != nil || &order == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) HandleDeleteOrder(w http.ResponseWriter, r *http.Request, orderID string) {
	order, err := h.orderService.DeleteOrder(orderID)
	if err != nil || &order == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
