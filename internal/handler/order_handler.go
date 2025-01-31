package handler

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"hot-coffee/utils"
	"net/http"
)

type OrderHandlerInterface interface {
	HandleCreateOrder(w http.ResponseWriter, r *http.Request)
	HandleGetAllOrders(w http.ResponseWriter, r *http.Request)
	HandleGetOrderById(w http.ResponseWriter, r *http.Request, orderID string)
	HandleDeleteOrder(w http.ResponseWriter, r *http.Request, orderID string)
	HandleUpdateOrder(w http.ResponseWriter, r *http.Request, orderID string)
	HandleCloseOrder(w http.Response, r *http.Request, orderID string)
}

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		utils.ErrorInJSON(w, http.StatusBadRequest, fmt.Errorf("invalid JSON format: %v", err))
		return
	}

	if order, err := h.orderService.CreateOrder(newOrder); err != nil {
		utils.ErrorInJSON(w, 400, err)
		return
	} else {
		utils.ResponseInJSON(w, 201, order)
	}
}

func (h *OrderHandler) HandleGetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
	}

	utils.ResponseInJSON(w, 200, orders)
}

func (h *OrderHandler) HandleGetOrderById(w http.ResponseWriter, r *http.Request, orderID string) {
	order, err := h.orderService.GetOrderByID(orderID)
	if err != nil || &order == nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	}
	utils.ResponseInJSON(w, 200, order)
}

func (h *OrderHandler) HandleDeleteOrder(w http.ResponseWriter, r *http.Request, orderID string) {
	order, err := h.orderService.DeleteOrder(orderID)
	if err != nil || &order == nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	}
	utils.ResponseInJSON(w, 204, order)
}

func (h *OrderHandler) HandleUpdateOrder(w http.ResponseWriter, r *http.Request, orderID string) {
	var changeOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&changeOrder); err != nil {
		utils.ErrorInJSON(w, http.StatusBadRequest, err)
		return
	}

	if order, err := h.orderService.UpdateOrder(orderID, changeOrder); err != nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	} else {
		utils.ResponseInJSON(w, 200, order)
	}
}

func (h *OrderHandler) HandleCloseOrder(w http.ResponseWriter, r *http.Request, orderID string) {
	if order, err := h.orderService.CloseOrder(orderID); err != nil {
		utils.ErrorInJSON(w, http.StatusNotFound, err)
		return
	} else {
		utils.ResponseInJSON(w, 200, order)
	}
}
