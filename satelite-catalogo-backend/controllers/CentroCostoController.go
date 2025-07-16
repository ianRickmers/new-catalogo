package controllers

import (
	"catalogo-backend/models"
	"catalogo-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCentroCosto godoc
// @Summary      Create centro de costo
// @Description  Creates a new Centro de Costo
// @Tags         cc
// @Accept       json
// @Produce      json
// @Param        payload  body      models.CC  true  "Centro de Costo info"
// @Success      201      {object} map[string]string
// @Failure      400      {object} map[string]interface{}
// @Router       /cc/ [post]
func CreateCentroCosto(ctx *gin.Context) {
	var cc models.CC
	if err := ctx.ShouldBindJSON(&cc); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar los datos del centro de costo"})
		return
	}

	id, err := services.NewCentroCostoService().CreateCC(&cc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el centro de costo"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id.Hex()})
}

// GetCentroCostoByID godoc
// @Summary      Get centro de costo by ID
// @Description  Returns a Centro de Costo by its ID
// @Tags         cc
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Centro de Costo ID"
// @Success      200  {object} models.CC
// @Failure      400  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Router       /cc/{id} [get]
func GetCentroCostoByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de centro de costo requerido"})
		return
	}

	cc, err := services.NewCentroCostoService().GetCCByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el centro de costo"})
		return
	}

	if cc == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Centro de costo no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, cc)
}

// UpdateCentroCosto godoc
// @Summary      Update centro de costo
// @Description  Updates a Centro de Costo by ID
// @Tags         cc
// @Accept       json
// @Produce      json
// @Param        id      path      string                true  "Centro de Costo ID"
// @Param        payload body      map[string]interface{} true  "Update data"
// @Success      200     {object} map[string]string
// @Failure      400     {object} map[string]interface{}
// @Router       /cc/{id} [put]
func UpdateCentroCosto(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de centro de costo requerido"})
		return
	}

	var updateData map[string]interface{}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar los datos de actualización"})
		return
	}

	err := services.NewCentroCostoService().UpdateCC(id, updateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el centro de costo"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Centro de costo actualizado correctamente"})
}

// DeleteCentroCosto godoc
// @Summary      Delete centro de costo
// @Description  Deletes a Centro de Costo by ID
// @Tags         cc
// @Produce      json
// @Param        id   path      string  true  "Centro de Costo ID"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]interface{}
// @Router       /cc/{id} [delete]
func DeleteCentroCosto(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de centro de costo requerido"})
		return
	}

	err := services.NewCentroCostoService().DeleteCC(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el centro de costo"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Centro de costo eliminado correctamente"})
}

// GetAllCentroCostos godoc
// @Summary      List centros de costo
// @Description  Returns all Centros de Costo
// @Tags         cc
// @Produce      json
// @Success      200  {array}  models.CC
// @Failure      500  {object} map[string]interface{}
// @Router       /cc/ [get]
func GetAllCentroCostos(ctx *gin.Context) {
	ccList, err := services.NewCentroCostoService().GetAllCC()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los centros de costo"})
		return
	}

	ctx.JSON(http.StatusOK, ccList)
}
