package services

import (
	"catalogo-backend/models"
	"catalogo-backend/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var centroCostoService *CentroCostoService

type CentroCostoService struct {
	repo *repositories.CentroCostoRepository
}

func NewCentroCostoService() *CentroCostoService {
	if centroCostoService == nil {
		centroCostoService = &CentroCostoService{
			repo: repositories.NewCentroCostoRepository(),
		}
	}
	return centroCostoService
}
func (s *CentroCostoService) CreateCC(cc *models.CC) (primitive.ObjectID, error) {
	cc.ID = primitive.NewObjectID()

	id, err := s.repo.InsertOne(cc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if cc.Jefe != primitive.NilObjectID {
		userRepo := repositories.NewUserRepository()
		err := userRepo.AddCCToUser(cc.Jefe, cc.ID)
		if err != nil {
			return id, err
		}
	}

	return id, nil
}

func (s *CentroCostoService) GetCCIDsByJefe(jefeID primitive.ObjectID) ([]primitive.ObjectID, error) {
	centros, err := s.repo.FindByJefe(jefeID)
	if err != nil {
		return nil, err
	}
	ids := make([]primitive.ObjectID, 0, len(centros))
	for _, cc := range centros {
		ids = append(ids, cc.ID)
	}
	return ids, nil
}

func (s *CentroCostoService) GetCCByID(id string) (*models.CC, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	cc, err := s.repo.FindByID(objID)
	if err != nil {
		return nil, err
	}
	if cc == nil {
		return nil, nil // No se encontró el centro de costo
	}
	return cc, nil
}

func (s *CentroCostoService) UpdateCC(id string, update bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	err = s.repo.UpdateOne(filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *CentroCostoService) DeleteCC(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	err = s.repo.DeleteOne(filter)
	if err != nil {
		return err
	}
	return nil
}
func (s *CentroCostoService) GetAllCC() ([]*models.CC, error) {
	centrosCosto, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return centrosCosto, nil
}
func (s *CentroCostoService) GetByNumeroCC(numeroCC int) (*models.CC, error) {
	cc, err := s.repo.FindByNumero(numeroCC)
	if err != nil {
		return nil, err
	}
	if cc == nil {
		return nil, nil // No se encontró el centro de costo
	}
	return cc, nil
}
