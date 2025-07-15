package services

import (
	"catalogo-backend/models"
	"catalogo-backend/repositories"
	"sync"
	"time"
)

var (
	logService *LogService
	onceLog    sync.Once
)

type LogService struct {
	repo *repositories.LogRepository
}

func getLogService() *LogService {
	onceLog.Do(func() {
		logService = &LogService{
			repo: repositories.NewLogRepository(),
		}
	})
	return logService
}

// createLogFromSolicitud crea un log a partir de una solicitud
func CreateLogFromSolicitud(solicitud *models.Solicitud) (string, error) {
	logEntry := &models.RequestLog{
		RequestID:     solicitud.ID,
		Timestamp:     time.Now(),
		EventType:     "create",
		Description:   "creación de solicitud",
		PreviousState: nil,       // No hay estado previo al crear una solicitud
		NewState:      solicitud, // El nuevo estado es la solicitud actual
		UserID:        solicitud.Solicitante,
	}

	id, err := getLogService().CreateLog(logEntry)
	if err != nil {
		return "error al crear el log de la solicitud", err
	}
	return id, nil
}

// createLogFromUpdate crea un log a partir de una actualización de solicitud
func CreateLogFromUpdate(solicitud *models.Solicitud, previousState *models.Solicitud) (string, error) {
	logEntry := &models.RequestLog{
		RequestID:     solicitud.ID,
		Timestamp:     time.Now(),
		EventType:     "update",
		Description:   "Actualización de la solicitud",
		PreviousState: previousState,       // Estado previo antes de la actualización
		NewState:      solicitud,           // El nuevo estado es la solicitud actualizada
		UserID:        solicitud.Aprobador, // se asocia al aprobador que realizó la actualización
	}

	id, err := getLogService().CreateLog(logEntry)
	if err != nil {
		return "error al crear el log de actualización de la solicitud", err
	}
	return id, nil
}

// Métodos CRUD clásicos
func (s *LogService) CreateLog(log *models.RequestLog) (string, error) {
	return s.repo.InsertOne(log)
}

func (s *LogService) GetAllLogs() ([]*models.RequestLog, error) {
	return s.repo.FindAll()
}

func (s *LogService) GetLogByID(id string) (*models.RequestLog, error) {
	return s.repo.FindByID(id)
}

func (s *LogService) UpdateLog(id string, log *models.RequestLog) error {
	return s.repo.UpdateOne(id, log)
}

func (s *LogService) DeleteByID(id string) error {
	return s.repo.DeleteByID(id)
}
