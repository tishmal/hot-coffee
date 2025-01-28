package dal

type InventoryRepositoryInterface interface {
}

type InventoryRepositoryJSON struct {
	filePath string
}

func NewInventoryRepositoryJSON(filepath string) InventoryRepositoryJSON {
	return InventoryRepositoryJSON{filePath: filepath}
}
