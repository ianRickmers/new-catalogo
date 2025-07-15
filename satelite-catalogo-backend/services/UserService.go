// services/UserService.go
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

// Repositorio global para usuarios
// se usa sync.Once para que el repositorio se inicialize solo una vez, después de que InitMongo() haya sido ejecutado
var (
	userRepo *repositories.UserRepository
	onceUser sync.Once
)

// Esto solo se ejecutará una vez, sin importar cuántas veces se llame gracias a sync.Once
func getUserRepo() *repositories.UserRepository {
	onceUser.Do(func() {
		userRepo = repositories.NewUserRepository()
	})
	return userRepo
}

func CreateUserService(newUser *models.User) (*models.User, error) {
	//Verificar si el usuario ya existe
	existingUser, err := getUserRepo().FindOne(bson.M{"email": newUser.Email})
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("el usuario ya existe con el email: %s", newUser.Email)
	}
	idUser, err := getUserRepo().InsertOne(newUser)
	if err != nil {
		return nil, err
	}
	newUser.ID = idUser
	return newUser, nil
}

func GetUserByUsernameService(username string) (*models.User, error) {
	utils.Debug("Iniciar usuario")

	filter := bson.M{"username": username}
	user, err := getUserRepo().FindOne(filter)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByIdService(id string) (models.User, error) {
	utils.Debug("Buscar usuario por id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, fmt.Errorf("formato de ID inválido: %s", id)
	}

	filter := bson.M{"_id": objID}
	user, err := getUserRepo().FindOne(filter)
	if err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func GetUserByEmailService(email string) (models.User, error) {
	utils.Debug("Buscar usuario por email")
	filter := bson.M{"email": email}
	user, err := getUserRepo().FindOne(filter)
	if err != nil {
		return models.User{}, err
	}

	if user == nil {
		return models.User{}, fmt.Errorf("usuario no encontrado con el email: %s", email)
	}

	return *user, nil
}

func GetAllUsersService() ([]*models.User, error) {
	utils.Debug("Get all users")

	users, err := getUserRepo().FindAllFiltered(bson.M{})
	if err != nil {

		return nil, err
	}
	return users, nil
}

func UpdateUserService(updatedUser models.User, userEmail string) (models.User, error) {
	utils.Debug("Update user")

	err := getUserRepo().UpdateOne(bson.M{"email": userEmail}, bson.M{"$set": updatedUser})
	if err != nil {
		return models.User{}, err
	}
	return updatedUser, nil
}

func DeleteUserService(userEmail string) error {
	utils.Debug("Delete user")

	err := getUserRepo().DeleteOne(userEmail)
	if err != nil {
		return err
	}
	return nil
}

func GetUsersByRoleService(role string) ([]*models.User, error) {
	utils.Debug("Get users by role")

	users, err := getUserRepo().FindAllFiltered(bson.M{"role": role})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func CheckRUTService(rut string) (bool, error) {
	utils.Debug("Check RUT")

	filter := bson.M{"rut": rut}
	user, err := getUserRepo().FindOne(filter)
	if err != nil {
		return false, nil
	}
	if user != nil {
		return true, nil
	}
	return false, nil
}

func CheckEmailService(email string) (bool, error) {
	utils.Debug("Check email")

	filter := bson.M{"email": email}
	user, err := getUserRepo().FindOne(filter)
	if err != nil {
		return false, nil
	}
	if user != nil {
		return true, nil
	}
	return false, nil
}

func GetUsersByCCService(ccIDs []primitive.ObjectID) ([]*models.User, error) {

	users, err := getUserRepo().FindByCC(ccIDs)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("no se encontraron usuarios para los IDs de centro de costo proporcionados")
	}
	return users, nil
}
