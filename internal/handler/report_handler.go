package handler

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/service"
	"hot-coffee/utils"
	"net/http"
)

type ReportHandlerInterface interface {
	HandleGetTotalSales(w http.ResponseWriter, r *http.Request)
	HandleGetPopulatItem(w http.ResponseWriter, r *http.Request)
}

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) ReportHandler {
	return ReportHandler{reportService: reportService}
}

func (h *ReportHandler) HandleGetTotalSales(w http.ResponseWriter, r *http.Request) {
	totalSales, err := h.reportService.GetTotalSales()
	if err != nil {
		http.Error(w, "Failed to retrieve total sales", http.StatusInternalServerError)
		fmt.Println("Error fetching total sales:", err)
		return
	}

	response := map[string]float64{"total sales": totalSales}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *ReportHandler) HandleGetPopulatItem(w http.ResponseWriter, r *http.Request) {
	popularItem, err := h.reportService.GetPopularItems()
	if err != nil {
	}
	utils.ResponseInJSON(w, 200, popularItem)
}
