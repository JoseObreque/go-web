package main

import (
	"github.com/JoseObreque/go-web/cmd/server/handler"
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
	router := gin.Default()

	// Ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Products endpoints
	productGroup := router.Group("/products")
	{
		productGroup.GET("/all", productHandler.GetAll())
		productGroup.GET("/:id", productHandler.GetById())
		productGroup.GET("/search", productHandler.GetByPriceGt())
		productGroup.POST("/new", productHandler.Create())
		productGroup.PUT("/:id", productHandler.FullUpdate())
		productGroup.PATCH("/:id", productHandler.PartialUpdate())
		productGroup.DELETE("/:id", productHandler.Delete())
	}

	// Start server
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
