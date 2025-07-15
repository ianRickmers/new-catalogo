package repositories

import (
	"context"
	"log"

	"catalogo-backend/database"
	"catalogo-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var logRepo *LogRepository

type LogRepository struct {
	collection *mongo.Collection
}

// Constructor
func NewLogRepository() *LogRepository {
	if database.Client == nil {
		log.Fatal("MongoDB client not initialized. Call InitMongo() first.")
	}

	if logRepo == nil {
		log.Println("Inicializando LogRepository")
		db := database.GetDatabase()
		collection := db.Collection("logs")
		logRepo = &LogRepository{collection: collection}
	}

	return logRepo
}

// Obtener todos los logs
func (repo *LogRepository) FindAll() ([]*models.RequestLog, error) {
	ctx := context.Background()
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*models.RequestLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// Insertar un log
func (repo *LogRepository) InsertOne(logEntry *models.RequestLog) (string, error) {
	ctx := context.Background()
	result, err := repo.collection.InsertOne(ctx, logEntry)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// Buscar un log por ID
func (repo *LogRepository) FindByID(id string) (*models.RequestLog, error) {
	ctx := context.Background()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var logEntry models.RequestLog
	err = repo.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&logEntry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &logEntry, nil
}

// Eliminar un log por ID
func (repo *LogRepository) DeleteByID(id string) error {
	ctx := context.Background()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = repo.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// Actualizar un log
func (repo *LogRepository) UpdateOne(id string, logEntry *models.RequestLog) error {
	ctx := context.Background()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": logEntry,
	}

	result, err := repo.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
