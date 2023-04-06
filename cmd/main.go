package main

import (
	docs "github.com/JoseObreque/go-web/cmd/docs"
	"github.com/JoseObreque/go-web/cmd/server/handler"
	"github.com/JoseObreque/go-web/cmd/server/middleware"
	"github.com/JoseObreque/go-web/internal/product"
	"github.com/JoseObreque/go-web/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// @BasePath /api/v1

// @title MELI Bootcamp API
// @version 1.0
// @description This API handles MELI products data.
// @termsOfService https://developers.mercadolibre.cl/es_ar/terminos-y-condiciones

// @contact.name API Support
// @contact.url https://developers.mercadolibre.cl/es_ar/support
func main() {
	// Load environment variables
	err := godotenv.Load("./cmd/local.env")
	if err != nil {
		panic(err)
	}

	// Extract products data from the JSON file
	jsonStore := store.NewJsonStore("products.json")
	productList, err := jsonStore.GetAll()
	if err != nil {
		panic(err)
	}

	// New product handler initialization
	repository := product.NewRepository(productList)
	service := product.NewService(repository)
	productHandler := handler.NewProductHandler(service)

	// Create new router
	router := gin.New()
	router.Use(middleware.PanicLogger())
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Products endpoints
	generalGroup := router.Group("/api/v1")

	// Ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Panic endpoint
	router.GET("/panic", func(c *gin.Context) {
		panic("oh no!")
	})

	// Swagger documentation endpoint
	generalGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Products endpoints
	productGroup := generalGroup.Group("/products")
	{
		productGroup.GET("/all", productHandler.GetAll())
		productGroup.GET("/:id", productHandler.GetById())
		productGroup.GET("/search", productHandler.GetByPriceGt())
	}

	protectedProductGroup := generalGroup.Group("/products")
	protectedProductGroup.Use(middleware.TokenValidator())
	{
		protectedProductGroup.POST("/new", productHandler.Create())
		protectedProductGroup.PUT("/:id", productHandler.FullUpdate())
		protectedProductGroup.PATCH("/:id", productHandler.PartialUpdate())
		protectedProductGroup.DELETE("/:id", productHandler.Delete())
	}

	// Start server
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
