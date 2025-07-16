package controllers

import (
	"fmt"
	"net/http"
	"time"

	"catalogo-backend/middleware"
	"catalogo-backend/models"
	"catalogo-backend/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Se establecen los nombres de la colección que se traeran desde la base de datos
*/

// CreateUserRequest represents the payload to create a new user.
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Rut      string `json:"rut"`
	CC       []int  `json:"cc"`
}

// CreateUser godoc
// @Summary      Register new user
// @Description  Creates a new user after validating credentials
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        payload  body      CreateUserRequest  true  "User info"
// @Success      201      {object} models.User
// @Failure      400      {object} map[string]interface{}
// @Router       /user/ [post]
func CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar los datos del usuario"})
		return
	}

	// Verificar con login externo (USACH)
	response, err := middleware.RequestLogin(req.Username, req.Password)

	if response == nil || response.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al verificar el usuario"})
		return
	}

	newUser := &models.User{
		Username:  req.Username,
		Email:     req.Username + "@usach.cl",
		Rut:       req.Rut,
		Role:      []models.Role{models.USER},
		CC:        []primitive.ObjectID{},
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	createdUser, err := services.CreateUserService(newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

func GetUserById(ctx *gin.Context) {
	userID := ctx.Param("id")
	resultUser, err := services.GetUserByIdService(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener el usuario"})
		return
	}
	ctx.JSON(http.StatusCreated, resultUser)
}

func GetUserByEmail(ctx *gin.Context) {
	userEmail := ctx.Param("email")
	resultUser, err := services.GetUserByEmailService(userEmail)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener el usuario"})
		return
	}
	ctx.JSON(http.StatusCreated, resultUser)
}

func GetAllUsers(ctx *gin.Context) {
	resultUser, err := services.GetAllUsersService()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener los usuarios"})
		return
	}
	ctx.JSON(http.StatusCreated, resultUser)
}

func UpdateUser(ctx *gin.Context) {
	var updatedUser models.User
	if err := ctx.ShouldBindJSON(&updatedUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar los datos de usuario"})
		return
	}
	userEmail := ctx.Param("email")
	updatedUser, err := services.UpdateUserService(updatedUser, userEmail)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al actualizar el usuario"})
		return
	}
	ctx.JSON(http.StatusCreated, "Usuario actualizado")
}

func DeleteUser(ctx *gin.Context) {
	userEmail := ctx.Param("email")
	err := services.DeleteUserService(userEmail)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Usuario eliminado")
}

func CheckRUT(ctx *gin.Context) {
	rut := ctx.Param("rut")
	resultUser, err := services.CheckRUTService(rut)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener el usuario"})
		return
	}
	ctx.JSON(http.StatusCreated, resultUser)
}

func CheckEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	resultUser, err := services.CheckEmailService(email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener el usuario"})
		return
	}
	ctx.JSON(http.StatusCreated, resultUser)
}

func GetUsersByCC(ctx *gin.Context) {
	var ccIDs []string
	if err := ctx.ShouldBindJSON(&ccIDs); err != nil {
		fmt.Println("Error real del BindJSON:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var objectIDs []primitive.ObjectID
	for _, idStr := range ccIDs {
		oid, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de CC inválido: " + idStr})
			return
		}
		objectIDs = append(objectIDs, oid)
	}

	resultUser, err := services.GetUsersByCCService(objectIDs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener los usuarios por CC"})
		return
	}
	ctx.JSON(http.StatusOK, resultUser)
}
