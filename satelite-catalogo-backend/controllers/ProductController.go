package controllers

import (
	"catalogo-backend/models"
	"catalogo-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateProduct godoc
// @Summary      Create product
// @Description  Creates a new product with the provided information
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        payload  body      models.Product  true  "Product info"
// @Success      201      {object} models.Product
// @Failure      400      {object} map[string]interface{}
// @Router       /product/ [post]
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

// GetAllProducts godoc
// @Summary      List products
// @Description  Returns all products
// @Tags         products
// @Produce      json
// @Success      200  {array} models.Product
// @Failure      500  {object} map[string]interface{}
// @Router       /product/ [get]
func GetAllProducts(ctx *gin.Context) {
	products, err := services.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary      Get product by ID
// @Description  Returns a product by its ID
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object} models.Product
// @Failure      404  {object} map[string]interface{}
// @Router       /product/{id} [get]

func GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")
	product, err := services.GetProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary      Update product
// @Description  Updates an existing product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Product ID"
// @Param        payload body      models.Product  true  "Product info"
// @Success      200     {object} models.Product
// @Failure      400     {object} map[string]interface{}
// @Router       /product/{id} [put]

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

// DeleteProduct godoc
// @Summary      Delete product
// @Description  Deletes a product by ID
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      204  {string} string "No Content"
// @Failure      500  {object} map[string]interface{}
// @Router       /product/{id} [delete]

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

// GetProductsPaginated godoc
// @Summary      List products paginated
// @Description  Returns paginated products
// @Tags         products
// @Produce      json
// @Param        page      query     int  false  "Page number"
// @Param        pageSize  query     int  false  "Page size"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /product/paginated [get]
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

// GetProductsFiltradasPaginated godoc
// @Summary      List filtered products paginated
// @Description  Returns products filtered by category, id or description
// @Tags         products
// @Produce      json
// @Param        page       query     int     false  "Page number"
// @Param        pageSize   query     int     false  "Page size"
// @Param        categoria  query     string  false  "Categoria"
// @Param        id_product query     string  false  "Product ID"
// @Param        descripcion query    string  false  "Descripcion"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /product/filtradas [get]

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
