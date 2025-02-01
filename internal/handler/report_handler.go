package handler

import (
	"fmt"
	"hot-coffee/internal/service"
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
	}
	fmt.Println(totalSales)
}

func (h *ReportHandler) HandleGetPopulatItem(w http.ResponseWriter, r *http.Request) {
}
