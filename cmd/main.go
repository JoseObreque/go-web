package main

import (
	"github.com/JoseObreque/go-web/cmd/server/handler"
	"github.com/JoseObreque/go-web/cmd/server/middleware"
	"github.com/JoseObreque/go-web/internal/product"
	"github.com/JoseObreque/go-web/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
)

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

	// Ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Panic endpoint
	router.GET("/panic", func(c *gin.Context) {
		panic("oh no!")
	})

	// Products endpoints
	productGroup := router.Group("/products")
	{
		productGroup.GET("/all", productHandler.GetAll())
		productGroup.GET("/:id", productHandler.GetById())
		productGroup.GET("/search", productHandler.GetByPriceGt())

	}

	protectedProductGroup := router.Group("/products")
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
