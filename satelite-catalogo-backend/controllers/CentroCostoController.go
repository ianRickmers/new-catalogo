package controllers

import (
	"catalogo-backend/models"
	"catalogo-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCentroCosto handles the creation of a new Centro de Costo
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
func GetAllCentroCostos(ctx *gin.Context) {
	ccList, err := services.NewCentroCostoService().GetAllCC()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los centros de costo"})
		return
	}

	ctx.JSON(http.StatusOK, ccList)
}
