package controllers

import (
	"catalogo-backend/models"
	"catalogo-backend/services"
	"catalogo-backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateSolicitud godoc
// @Summary      Create solicitud
// @Description  Creates a new solicitud with optional files
// @Tags         solicitudes
// @Accept       multipart/form-data
// @Produce      json
// @Param        solicitud  formData  string  true  "Solicitud JSON"
// @Param        archivos   formData  file    false "Attached files"
// @Success      201  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Router       /solicitud/ [post]
func CreateSolicitud(ctx *gin.Context) {
	// se envia como formData ya que recibe los archivos como multipart/form-data
	jsonStr := ctx.PostForm("solicitud")
	var solicitud models.Solicitud
	if err := json.Unmarshal([]byte(jsonStr), &solicitud); err != nil {
		utils.Debug("Error al deserializar solicitud:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida"})
		return
	}
	// le damos un ID único a la solicitud
	solicitud.ID = primitive.NewObjectID()

	form, _ := ctx.MultipartForm()
	files := form.File["archivos"]

	if len(files) > 0 {
		// la carpeta es generada con el ID de la solicitud en Hexadecimal como string
		rutas, err := utils.GuardarArchivos(files, solicitud.ID.Hex())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar archivos"})
			return
		}
		solicitud.Documents = rutas
	}
	// se crea la solicitud en la base de datos
	result, err := services.CreateSolicitudService(&solicitud)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error al crear la solicitud": err.Error()})
		return
	}
	// se crea el log de la solicitud
	resultado, err := services.CreateLogFromSolicitud(&solicitud)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear log de solicitud: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"id":         result.ID,
		"documentos": solicitud.Documents,
		"log":        resultado,
	})
}

// GetSolicitud godoc
// @Summary      Get solicitud by ID
// @Description  Returns a solicitud by its ID
// @Tags         solicitudes
// @Produce      json
// @Param        id   path      string  true  "Solicitud ID"
// @Success      200  {object} models.Solicitud
// @Failure      404  {object} map[string]interface{}
// @Router       /solicitud/{id} [get]
func GetSolicitud(ctx *gin.Context) {
	id := ctx.Param("id")

	solicitud, err := services.GetSolicitudByIDService(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if solicitud == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Solicitud no encontrada"})
		return
	}

	ctx.JSON(http.StatusOK, solicitud)
}

func GetAllSolicitudes(ctx *gin.Context) {
	solicitudes, err := services.GetAllSolicitudesService()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, solicitudes)
}

// UpdateSolicitud godoc
// @Summary      Update solicitud
// @Description  Updates an existing solicitud
// @Tags         solicitudes
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Solicitud ID"
// @Param        payload body      object  true  "Update data"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]interface{}
// @Router       /solicitud/{id} [put]
func UpdateSolicitud(ctx *gin.Context) {
	id := ctx.Param("id")

	var update bson.M
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos para actualización"})
		return
	}
	// antes de actualizar, obtenemos la solicitud actual para crear el log
	solicitudPrevia, err := services.GetSolicitudByIDService(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener solicitud para obtener estado previo (logs): " + err.Error()})
		return
	}
	if err := services.UpdateSolicitudService(id, update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error al actualizar la solicitud": err.Error()})
		return
	}
	// se crea el log de la solicitud
	solicitudPosterior, err := services.GetSolicitudByIDService(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener solicitud posterior a la actualización (logs): " + err.Error()})
	}
	// se crea el log de actualización
	_, err = services.CreateLogFromUpdate(solicitudPosterior, solicitudPrevia)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear log de actualización: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Solicitud actualizada correctamente"})
}

// DeleteSolicitud godoc
// @Summary      Delete solicitud
// @Description  Deletes a solicitud by ID
// @Tags         solicitudes
// @Produce      json
// @Param        id   path      string  true  "Solicitud ID"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]interface{}
// @Router       /solicitud/{id} [delete]
func DeleteSolicitud(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := services.DeleteSolicitudService(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Solicitud eliminada"})
}

// GetSolicitudesFiltradasPaginated godoc
// @Summary      List solicitudes filtered paginated
// @Description  Returns solicitudes using filters and pagination
// @Tags         solicitudes
// @Produce      json
// @Param        page       query int    false "Page"
// @Param        pageSize   query int    false "Page size"
// @Param        state      query string false "State"
// @Param        id         query string false "Solicitud ID"
// @Param        fechaInicio query string false "Fecha inicio"
// @Param        fechaFin   query string false "Fecha fin"
// @Param        ccs        query []string false "Centros de costo"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /solicitud/filtradas [get]
func GetSolicitudesFiltradasPaginated(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "50")
	state := ctx.Query("state")
	idStr := ctx.Query("id")
	fechaInicioStr := ctx.Query("fechaInicio")
	fechaFinStr := ctx.Query("fechaFin")
	// Obtener los centros de costo enviados desde el frontend
	ccIDs := ctx.QueryArray("ccs[]")
	if len(ccIDs) == 0 {
		ccIDs = ctx.QueryArray("ccs") // fallback
		// si no llegan centros de costo, no se muestra nada
		if len(ccIDs) == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"data":       []models.Solicitud{},
				"total":      0,
				"page":       1,
				"pageSize":   50,
				"totalPages": 0,
			})
			return
		}
	}
	fmt.Println("Centros de Costo recibidos:", ccIDs)

	// Paginación segura
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}

	filter := bson.M{}

	if state != "" {
		filter["state"] = state
	}

	if idStr != "" {
		objID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"data":       []models.Solicitud{},
				"total":      0,
				"page":       page,
				"pageSize":   pageSize,
				"totalPages": 0,
			})
			return
		}
		filter["_id"] = objID
	}
	// Si ambas fechas están presentes, intentamos parsear y agregar al filtro
	if fechaInicioStr != "" && fechaFinStr != "" {
		fechaInicio, err1 := time.Parse(time.RFC3339, fechaInicioStr)
		fechaFin, err2 := time.Parse(time.RFC3339, fechaFinStr)
		utils.Debug("Fecha Inicio:", fechaInicioStr, "Fecha Fin:", fechaFinStr)
		if err1 == nil && err2 == nil {
			fechaInicio = time.Date(fechaInicio.Year(), fechaInicio.Month(), fechaInicio.Day(), 0, 0, 0, 0, time.UTC)
			fechaFin = time.Date(fechaFin.Year(), fechaFin.Month(), fechaFin.Day(), 0, 0, 0, 0, time.UTC)
			// Agregar filtro por fechas (entre fechaInicio y fin del día de fechaFin)

			filter["fecha_solicitud"] = bson.M{
				"$gte": fechaInicio,
				"$lt":  fechaFin.AddDate(0, 0, 1), // Para incluir todo el día
			}
			utils.Debug("Filtro por fechas aplicado:", filter["fecha_solicitud"])
		}
	}
	// Aplicar filtro por centros de costo si se recibieron
	if len(ccIDs) > 0 {
		var objectIDs []primitive.ObjectID
		for _, idStr := range ccIDs {
			if oid, err := primitive.ObjectIDFromHex(idStr); err == nil {
				objectIDs = append(objectIDs, oid)
			}
		}
		if len(objectIDs) > 0 {
			filter["cc"] = bson.M{"$in": objectIDs}
		}
	}

	// Llamada al servicio
	solicitudes, total, err := services.GetSolicitudesFilteredPaginatedService(page, pageSize, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       solicitudes,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int((total + int64(pageSize) - 1) / int64(pageSize)),
	})
}

// GetSolicitudesPaginated godoc
// @Summary      List solicitudes paginated
// @Description  Returns solicitudes paginated
// @Tags         solicitudes
// @Produce      json
// @Param        page      query int false "Page"
// @Param        pageSize  query int false "Page size"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /solicitud/ [get]
func GetSolicitudesPaginated(ctx *gin.Context) {
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

	solicitudes, total, err := services.GetSolicitudesPaginatedService(page, pageSize, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       solicitudes,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int((total + int64(pageSize) - 1) / int64(pageSize)), // redondeo hacia arriba
	})
}

func GetSolicitudesByCCAndStatePaginated(ctx *gin.Context) {
	ccRaw, exists := ctx.Get("cc")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Centro de costo no disponible"})
		return
	}

	cc, ok := ccRaw.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer centro de costo"})
		return
	}

	state := ctx.DefaultQuery("state", "I")

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "50"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	solicitudes, total, err := services.GetSolicitudesByCCAndStatePaginatedService(cc, state, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       solicitudes,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int((total + int64(pageSize) - 1) / int64(pageSize)),
	})
}

// GetSolicitudesAprobarPaginated godoc
// @Summary      List solicitudes to approve
// @Description  Returns solicitudes for approval filtered by supervisor
// @Tags         solicitudes
// @Produce      json
// @Param        userId    query string true  "User ID"
// @Param        page      query int    false "Page"
// @Param        pageSize  query int    false "Page size"
// @Param        state     query string false "State"
// @Param        id        query string false "Solicitud ID"
// @Param        cc        query string false "Centro de costo"
// @Param        fechaInicio query string false "Fecha inicio"
// @Param        fechaFin query string false "Fecha fin"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /solicitud/aprobar [get]
func GetSolicitudesAprobarPaginated(ctx *gin.Context) {
	userIDStr := ctx.Query("userId")
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Falta el parámetro userId"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de userId inválido"})
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "50")
	state := ctx.Query("state")
	idStr := ctx.Query("id")

	ccStr := ctx.Query("cc")
	fechaInicioStr := ctx.Query("fechaInicio")
	fechaFinStr := ctx.Query("fechaFin")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}

	ccIDs, err := services.NewCentroCostoService().GetCCIDsByJefe(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener centros de costo"})
		return
	}
	if len(ccIDs) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.Solicitud{},
			"total":      0,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": 0,
		})
		return
	}

	filter := bson.M{
		"cc": bson.M{"$in": ccIDs},
	}

	if state != "" {
		filter["state"] = state
	}

	if idStr != "" {
		if objID, err := primitive.ObjectIDFromHex(idStr); err == nil {
			filter["_id"] = objID
		}
	}

	if ccStr != "" {
		if ccID, err := primitive.ObjectIDFromHex(ccStr); err == nil {
			found := false
			for _, id := range ccIDs {
				if id == ccID {
					found = true
					break
				}
			}
			if found {
				filter["cc"] = ccID
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"data":       []models.Solicitud{},
					"total":      0,
					"page":       page,
					"pageSize":   pageSize,
					"totalPages": 0,
				})
				return
			}
		}
	}
	// Si ambas fechas están presentes, intentamos parsear y agregar al filtro
	if fechaInicioStr != "" && fechaFinStr != "" {
		fechaInicio, err1 := time.Parse(time.RFC3339, fechaInicioStr)
		fechaFin, err2 := time.Parse(time.RFC3339, fechaFinStr)
		utils.Debug("Fecha Inicio:", fechaInicioStr, "Fecha Fin:", fechaFinStr)
		if err1 == nil && err2 == nil {
			fechaInicio = time.Date(fechaInicio.Year(), fechaInicio.Month(), fechaInicio.Day(), 0, 0, 0, 0, time.UTC)
			fechaFin = time.Date(fechaFin.Year(), fechaFin.Month(), fechaFin.Day(), 0, 0, 0, 0, time.UTC)
			// Agregar filtro por fechas (entre fechaInicio y fin del día de fechaFin)

			filter["fecha_solicitud"] = bson.M{
				"$gte": fechaInicio,
				"$lt":  fechaFin.AddDate(0, 0, 1), // Para incluir todo el día
			}
			utils.Debug("Filtro por fechas aplicado:", filter["fecha_solicitud"])
		}
	}
	solicitudes, total, err := services.GetSolicitudesFilteredPaginatedService(page, pageSize, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       solicitudes,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int((total + int64(pageSize) - 1) / int64(pageSize)),
	})
}
