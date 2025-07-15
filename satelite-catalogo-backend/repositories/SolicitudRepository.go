package repositories

import (
	"context"
	"errors"
	"log"

	"catalogo-backend/database"
	"catalogo-backend/models"
	"catalogo-backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var solicitudRepo *SolicitudRepository

type SolicitudRepository struct {
	collection *mongo.Collection
}

func NewSolicitudRepository() *SolicitudRepository {
	if database.Client == nil {
		log.Fatal("MongoDB client not initialized. Call InitMongo() first.")
	}

	if solicitudRepo == nil {
		log.Println("Inicializando SolicitudRepository")
		db := database.GetDatabase()
		collection := db.Collection("solicitudes")
		solicitudRepo = &SolicitudRepository{collection: collection}
	}
	return solicitudRepo
}

func (repo *SolicitudRepository) InsertOne(solicitud *models.Solicitud) error {
	_, err := repo.collection.InsertOne(context.Background(), solicitud)
	return err
}

func (repo *SolicitudRepository) FindOne(filter bson.M) (*models.Solicitud, error) {
	var solicitud models.Solicitud
	err := repo.collection.FindOne(context.Background(), filter).Decode(&solicitud)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &solicitud, nil
}

func (repo *SolicitudRepository) UpdateOne(filter, update bson.M) error {
	result, err := repo.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no se encontró la solicitud")
	}
	return nil
}

func (repo *SolicitudRepository) FindAll() ([]*models.Solicitud, error) {
	var solicitudes []*models.Solicitud
	cursor, err := repo.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.Background(), &solicitudes); err != nil {
		return nil, err
	}
	return solicitudes, nil
}

// index solicitudes
func (repo *SolicitudRepository) FindFilteredPaginated(page, pageSize int, filter bson.M) ([]*models.Solicitud, int64, error) {
	var solicitudes []*models.Solicitud

	total, err := repo.collection.CountDocuments(context.Background(), filter)
	if err != nil {

		return nil, 0, err
	}

	opts := options.Find()
	opts.SetSkip(int64((page - 1) * pageSize))
	opts.SetLimit(int64(pageSize))
	opts.SetSort(bson.D{{Key: "fecha_solicitud", Value: -1}})

	cursor, err := repo.collection.Find(context.Background(), filter, opts)
	if err != nil {
		utils.Debug("Error al buscar solicitudes:", err)
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &solicitudes); err != nil {
		return nil, 0, err
	}

	return solicitudes, total, nil
}

func (repo *SolicitudRepository) DeleteOneByID(id primitive.ObjectID) error {
	result, err := repo.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no se encontró la solicitud")
	}
	return nil
}

func (repo *SolicitudRepository) FindAllPaginated(page, pageSize int, filter bson.M) ([]*models.Solicitud, int64, error) {
	var solicitudes []*models.Solicitud

	// Obtener el número total de documentos que cumplen con el filtro
	total, err := repo.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		utils.Debug("Error al contar documentos:", err)
		return nil, 0, err
	}

	// Configurar las opciones para paginar
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	opts := options.Find()
	opts.SetSkip(skip)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Orden descendente por fecha

	// Ejecutar la consulta
	cursor, err := repo.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &solicitudes); err != nil {
		return nil, 0, err
	}

	return solicitudes, total, nil
}

func (repo *SolicitudRepository) FindByCCAndStatePaginated(cc int, state string, page, pageSize int) ([]*models.Solicitud, int64, error) {
	filter := bson.M{
		"cc":    cc,
		"state": state,
	}

	totalCount, err := repo.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "fecha_actualizacion", Value: -1}})

	cursor, err := repo.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var solicitudes []*models.Solicitud
	if err := cursor.All(context.Background(), &solicitudes); err != nil {
		return nil, 0, err
	}

	return solicitudes, totalCount, nil
}
