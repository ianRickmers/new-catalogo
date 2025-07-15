package services

import (
	"catalogo-backend/models"
	"catalogo-backend/repositories"
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	productRepo *repositories.ProductRepository
	onceProduct sync.Once
)

func getProductRepo() *repositories.ProductRepository {
	onceProduct.Do(func() {
		productRepo = repositories.NewProductRepository()
	})
	return productRepo
}

// Create - Crear un nuevo producto
func CreateProduct(product models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return getProductRepo().Create(ctx, product)
}

// GetAll - Obtener todos los productos
func GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return getProductRepo().FindAll(ctx)
}

// GetByID - Obtener un producto por ID
func GetProductByID(id string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return getProductRepo().FindByID(ctx, id)
}

// Update - Actualizar un producto
func UpdateProduct(id string, product models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return getProductRepo().Update(ctx, id, product)
}

// Delete - Eliminar un producto
func DeleteProduct(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return getProductRepo().Delete(ctx, id)
}

// Search - Buscar productos
func SearchProduct(query string) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return getProductRepo().Search(ctx, query)
}

type PaginatedResult struct {
	Products     []models.Product `json:"products"`
	TotalRecords int64            `json:"totalRecords"`
	Page         int              `json:"page"`
	PageSize     int              `json:"pageSize"`
}

// GetAllPaginated - Obtener productos con paginación y filtros
// func GetAllPaginated(page, pageSize int, searchParams map[string]string) (*PaginatedResult, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Construir la query basada en los parámetros de búsqueda
// 	query := getProductRepo().BuildSearchQuery(searchParams)

// 	// Obtener los productos paginados
// 	products, totalRecords, err := getProductRepo().FindAllPaginated(ctx, page, pageSize, query)
// 	if err != nil {
// 		return nil, err
// 	}

//		return &PaginatedResult{
//			Products:     products,
//			TotalRecords: totalRecords,
//			Page:         page,
//			PageSize:     pageSize,
//		}, nil
//	}
func GetProductsPaginatedService(page, pageSize int, filter bson.M) ([]*models.Product, int64, error) {
	return getProductRepo().FindAllPaginated(page, pageSize, filter)
}
