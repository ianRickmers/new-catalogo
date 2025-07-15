package controllers

import (
	"catalogo-backend/models"
	"catalogo-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Create godoc
// Crear un nuevo producto
// Crea un nuevo producto con la información proporcionada

func CreateProduct(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateProduct(product); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

// GetAll godoc
//Obtener todos los productos
//Obtiene una lista de todos los productos

func GetAllProducts(ctx *gin.Context) {
	products, err := services.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GetByID godoc
// Obtener un producto por ID
//Obtiene un producto específico por su ID

func GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")
	product, err := services.GetProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// Update godoc
// Actualizar un producto
// Actualiza un producto existente por su ID

func UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateProduct(id, product); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// Delete godoc
// Eliminar un producto
// Elimina un producto por su ID

func DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := services.DeleteProduct(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetAllPaginated godoc
// Obtener productos paginados
// Obtiene una lista paginada de productos con filtros opcionales

// func GetAllProductsPaginated(ctx *gin.Context) {
// 	// Obtener parámetros de paginación
// 	pageStr := ctx.DefaultQuery("page", "1")
// 	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil || page < 1 {
// 		page = 1
// 	}

// 	pageSize, err := strconv.Atoi(pageSizeStr)
// 	if err != nil || pageSize < 1 {
// 		pageSize = 10
// 	}

// 	// Recopilar parámetros de búsqueda
// 	searchParams := map[string]string{
// 		"search":         ctx.Query("search"),
// 		"licitacion":     ctx.Query("licitacion"),
// 		"region":         ctx.Query("region"),
// 		"marca":          ctx.Query("marca"),
// 		"convenio_marco": ctx.Query("convenio_marco"),
// 	}

// 	// Obtener resultados paginados
// 	result, err := services.GetAllPaginated(page, pageSize, searchParams)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, result)
// }

func GetProductsPaginated(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "50")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 50
	}

	// filtrar por estado, si se proporciona
	filter := bson.M{}

	productos, total, err := services.GetProductsPaginatedService(page, pageSize, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       productos,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int((total + int64(pageSize) - 1) / int64(pageSize)), // redondeo hacia arriba
	})
}

func GetProductsFiltradasPaginated(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "50")
	categoria := ctx.Query("categoria")
	idStr := ctx.Query("id_product")
	descripcion := ctx.Query("descripcion")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 50
	}

	filter := bson.M{}

	if categoria != "" {
		filter["categoria"] = bson.M{"$regex": categoria, "$options": "i"}
	}

	if idStr != "" {
		filter["id_product"] = bson.M{"$regex": idStr, "$options": "i"}
	}

	if descripcion != "" {
		filter["descripcion"] = bson.M{"$regex": descripcion, "$options": "i"}
	}
	productos, total, err := services.GetProductsPaginatedService(page, pageSize, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       productos,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int((total + int64(pageSize) - 1) / int64(pageSize)),
	})
}
