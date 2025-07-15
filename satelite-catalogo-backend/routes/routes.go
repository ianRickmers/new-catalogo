package routes

import (
	"catalogo-backend/controllers"
	"catalogo-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Expone los archivos estáticos de uploads solo para usuarios autenticados
	archivosGroup := router.Group("/archivos")
	archivosGroup.Use(middleware.LoadJWTAuth().MiddlewareFunc())
	archivosGroup.GET("/*filepath", controllers.ServeArchivo)
	// User routes
	userGroup := router.Group("/user")
	userGroup.Use(middleware.LoadJWTAuth().MiddlewareFunc())
	{
		userGroup.POST("/", controllers.CreateUser)
		userGroup.GET("/:id", controllers.GetUserById)
		userGroup.GET("/email/:email", controllers.GetUserByEmail)
		userGroup.PUT("/:id", controllers.UpdateUser)
		userGroup.GET("/", controllers.GetAllUsers)
		userGroup.DELETE("/:id", controllers.DeleteUser)
		userGroup.POST("/by-cc", controllers.GetUsersByCC)
	}

	// Auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", middleware.LoadJWTAuth().LoginHandler)
		authGroup.POST("/refresh_token", middleware.LoadJWTAuth().RefreshHandler)
		authGroup.POST("/logout", middleware.LoadJWTAuth().LogoutHandler)
	}

	// Solicitud routes
	solicitudGroup := router.Group("/solicitud")
	solicitudGroup.Use(middleware.LoadJWTAuth().MiddlewareFunc())
	{
		solicitudGroup.GET("/filtradas", controllers.GetSolicitudesFiltradasPaginated)
		solicitudGroup.POST("/", controllers.CreateSolicitud)
		solicitudGroup.GET("/", controllers.GetSolicitudesPaginated)
		solicitudGroup.PUT("/:id", controllers.UpdateSolicitud)
		solicitudGroup.GET("/:id", controllers.GetSolicitud)
		solicitudGroup.DELETE("/:id", controllers.DeleteSolicitud)
		solicitudGroup.GET("/aprobar", controllers.GetSolicitudesAprobarPaginated)
	}
	// Centro de Costo routes
	ccGroup := router.Group("/cc")
	ccGroup.Use(middleware.LoadJWTAuth().MiddlewareFunc())
	{
		ccGroup.POST("/", controllers.CreateCentroCosto)
		ccGroup.GET("/:id", controllers.GetCentroCostoByID)
		ccGroup.PUT("/:id", controllers.UpdateCentroCosto)
		ccGroup.GET("/", controllers.GetAllCentroCostos)
		ccGroup.DELETE("/:id", controllers.DeleteCentroCosto)
	}

	products := router.Group("/product")
	products.Use(middleware.LoadJWTAuth().MiddlewareFunc())
	{
		products.POST("/", controllers.CreateProduct)
		products.GET("/", controllers.GetAllProducts)
		products.GET("/paginated", controllers.GetProductsPaginated)
		products.GET("/filtradas", controllers.GetProductsFiltradasPaginated)
		products.GET("/:id", controllers.GetProductByID)
		products.PUT("/:id", controllers.UpdateProduct)
		products.DELETE("/:id", controllers.DeleteProduct)
	}
}
