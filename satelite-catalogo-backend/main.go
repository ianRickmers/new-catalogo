// Package main bootstraps the Gin server, loads environment variables and
// establishes the database connection.
//
// @title Catalogo API
// @version 1.0
// @description This is the API documentation for the Catalogo service.
// @host localhost:8080
// @BasePath /
package main

import (
	docs "catalogo-backend/docs"
	"context"
	"log"
	"time"

	"catalogo-backend/database"
	"catalogo-backend/middleware"
	"catalogo-backend/routes"
	"catalogo-backend/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Cargar variables de entorno
	utils.LoadEnv()

	// Inicializar conexión Mongo
	database.InitMongo()

	// Desconectar al final
	defer func() {
		if err := database.Client.Disconnect(ctx); err != nil {
			log.Println("Error al desconectar MongoDB:", err)
		}
	}()

	r := gin.Default()

	r.Use(middleware.CorsMiddleware())

	docs.SwaggerInfo.BasePath = "/"

	// Registrar rutas
	routes.RegisterRoutes(r)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
