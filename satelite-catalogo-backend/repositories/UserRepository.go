// repositories/UserRepository.go
package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"catalogo-backend/database"
	"catalogo-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userRepo *UserRepository

type UserRepository struct {
	collection *mongo.Collection
}

// Constructor para UserRepository
func NewUserRepository() *UserRepository {
	if database.Client == nil {
		log.Fatal("MongoDB client not initialized. Call InitMongo() first.")
	}

	if userRepo == nil {

		log.Println("Inicializando UserRepository")
		db := database.GetDatabase()
		collection := db.Collection("users")
		userRepo = &UserRepository{collection: collection}
	}

	return userRepo
}

func (userRepo *UserRepository) InsertOne(user *models.User) (primitive.ObjectID, error) {
	inserted, err := userRepo.collection.InsertOne(context.Background(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return inserted.InsertedID.(primitive.ObjectID), nil
}

func (userRepo *UserRepository) FindOne(filter bson.M) (*models.User, error) {
	var user models.User
	err := userRepo.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (userRepo *UserRepository) UpdateOne(filter, update bson.M) error {
	updateResult, err := userRepo.collection.UpdateOne(context.Background(), filter, update)
	if updateResult.MatchedCount == 0 {
		return errors.New("no se encontró el usuario")
	}
	return err
}

func (userRepo *UserRepository) FindAll() ([]*models.User, error) {
	var users []*models.User
	cursor, err := userRepo.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (userRepo *UserRepository) FindAllFiltered(filter bson.M) ([]*models.User, error) {
	var users []*models.User
	cursor, err := userRepo.collection.Find(context.Background(), filter)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (userRepo *UserRepository) DeleteOne(email string) error {
	result, err := userRepo.collection.DeleteOne(context.Background(), bson.M{"email": email})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no se encontró el usuario")
	}

	return nil
}

func (r *UserRepository) FindByCC(ccIDs []primitive.ObjectID) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"cc": bson.M{
			"$in": ccIDs,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (userRepo *UserRepository) AddCCToUser(userID, ccID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$addToSet": bson.M{"cc": ccID}}

	_, err := userRepo.collection.UpdateOne(context.Background(), filter, update)
	return err
}
