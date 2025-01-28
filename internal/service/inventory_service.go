package service

import "hot-coffee/internal/dal"

type InventoryServiceInterface interface {
}

type InventoryService struct {
	Repository dal.InventoryRepositoryInterface
}

func NewInventoryService(repository dal.InventoryRepositoryInterface) InventoryService {
	return InventoryService{Repository: repository}
}
