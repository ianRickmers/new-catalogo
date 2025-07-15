package services

import (
	"fmt"
	"sync"

	"catalogo-backend/models"
	"catalogo-backend/repositories"
	"catalogo-backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	solicitudRepo *repositories.SolicitudRepository
	onceSolicitud sync.Once
)

func getSolicitudRepo() *repositories.SolicitudRepository {
	onceSolicitud.Do(func() {
		solicitudRepo = repositories.NewSolicitudRepository()
	})
	return solicitudRepo
}

func CreateSolicitudService(newSolicitud *models.Solicitud) (*models.Solicitud, error) {
	utils.Debug(fmt.Sprintf("Creando solicitud con ID %s", newSolicitud.ID.Hex()))

	err := getSolicitudRepo().InsertOne(newSolicitud)
	if err != nil {
		return nil, err
	}

	return newSolicitud, nil
}

func GetSolicitudByIDService(id string) (*models.Solicitud, error) {
	utils.Debug("Buscar solicitud por ID")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("formato de ID inválido: %s", id)
	}
	return getSolicitudRepo().FindOne(bson.M{"_id": objID})
}

func GetAllSolicitudesService() ([]*models.Solicitud, error) {
	utils.Debug("Obtener todas las solicitudes")

	return getSolicitudRepo().FindAll()
}

func UpdateSolicitudService(id string, update bson.M) error {
	utils.Debug("Actualizar solicitud")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("formato de ID inválido: %s", id)
	}
	return getSolicitudRepo().UpdateOne(bson.M{"_id": objID}, bson.M{"$set": update})
}

func DeleteSolicitudService(id string) error {
	utils.Debug("Eliminar solicitud")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("formato de ID inválido: %s", id)
	}
	return getSolicitudRepo().DeleteOneByID(objID)
}

// esta se usa en el index
func GetSolicitudesFilteredPaginatedService(page, pageSize int, filter bson.M) ([]*models.Solicitud, int64, error) {
	return getSolicitudRepo().FindFilteredPaginated(page, pageSize, filter)
}

func GetSolicitudesPaginatedService(page, pageSize int, filter bson.M) ([]*models.Solicitud, int64, error) {
	return getSolicitudRepo().FindAllPaginated(page, pageSize, filter)
}

func GetSolicitudesByCCAndStatePaginatedService(cc int, state string, page, pageSize int) ([]*models.Solicitud, int64, error) {
	return getSolicitudRepo().FindByCCAndStatePaginated(cc, state, page, pageSize)
}
