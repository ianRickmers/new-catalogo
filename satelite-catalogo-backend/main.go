// main.go
package main

import (
	"context"
	"log"
	"time"

	"catalogo-backend/database"
	"catalogo-backend/middleware"
	"catalogo-backend/routes"
	"catalogo-backend/utils"

	"github.com/gin-gonic/gin"
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

	// Registrar rutas
	routes.RegisterRoutes(r)

	r.Run(":8080")
}
