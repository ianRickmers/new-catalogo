package repositories

import (
	"context"
	"errors"
	"log"

	"catalogo-backend/database"
	"catalogo-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var centroCostoRepo *CentroCostoRepository

type CentroCostoRepository struct {
	collection *mongo.Collection
}

func NewCentroCostoRepository() *CentroCostoRepository {
	if database.Client == nil {
		log.Fatal("MongoDB client not initialized. Call InitMongo() first.")
	}

	if centroCostoRepo == nil {
		log.Println("Inicializando CentroCostoRepository")
		db := database.GetDatabase()
		collection := db.Collection("centros_costo")
		centroCostoRepo = &CentroCostoRepository{collection: collection}
	}
	return centroCostoRepo
}

func (repo *CentroCostoRepository) InsertOne(cc *models.CC) (primitive.ObjectID, error) {
	result, err := repo.collection.InsertOne(context.Background(), cc)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (repo *CentroCostoRepository) FindOne(filter bson.M) (*models.CC, error) {
	var cc models.CC
	err := repo.collection.FindOne(context.Background(), filter).Decode(&cc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &cc, nil
}
func (repo *CentroCostoRepository) UpdateOne(filter, update bson.M) error {
	result, err := repo.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (repo *CentroCostoRepository) DeleteOne(filter bson.M) error {
	result, err := repo.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (repo *CentroCostoRepository) FindAll() ([]*models.CC, error) {
	var centrosCosto []*models.CC
	cursor, err := repo.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &centrosCosto); err != nil {
		return nil, err
	}
	return centrosCosto, nil
}
func (repo *CentroCostoRepository) FindAllFiltered(filter bson.M) ([]*models.CC, error) {
	var centrosCosto []*models.CC
	cursor, err := repo.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &centrosCosto); err != nil {
		return nil, err
	}
	return centrosCosto, nil
}
func (repo *CentroCostoRepository) FindByID(id primitive.ObjectID) (*models.CC, error) {
	filter := bson.M{"_id": id}
	var cc models.CC
	err := repo.collection.FindOne(context.Background(), filter).Decode(&cc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &cc, nil
}
func (repo *CentroCostoRepository) FindByJefe(jefeID primitive.ObjectID) ([]*models.CC, error) {
	filter := bson.M{"jefe": jefeID}
	var centrosCosto []*models.CC
	cursor, err := repo.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &centrosCosto); err != nil {
		return nil, err
	}
	return centrosCosto, nil
}

func (repo *CentroCostoRepository) FindByNumero(numero int) (*models.CC, error) {
	filter := bson.M{"numero": numero}
	var cc models.CC
	err := repo.collection.FindOne(context.Background(), filter).Decode(&cc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &cc, nil
}
