package service

import "fmt"

type ReportService struct {
	menuService MenuService
	orderSevice OrderService
}

func NewReportService(menuService MenuService, orderService OrderService) ReportService {
	return ReportService{
		menuService: menuService,
		orderSevice: orderService,
	}
}

func (r *ReportService) GetTotalSales() (float64, error) {
	orders, err := r.orderSevice.GetAllOrders()
	if err != nil {
		return 0, fmt.Errorf("Failed to retrieve orders")
	}

	menuItems, err := r.menuService.GetAllMenuItems()
	if err != nil {
		return 0, fmt.Errorf("Failed to retrieve menu items")
	}

	totalSales := 0.0

	for _, order := range orders {
		for _, item := range order.Items {
			for _, menuItem := range menuItems {
				if item.ProductID == menuItem.ID {
					fmt.Println(item.Quantity, menuItem.Price)
					totalSales += float64(item.Quantity) * menuItem.Price
				}
			}
		}
	}

	return totalSales, nil
}

func GetPopularItems() {
}
