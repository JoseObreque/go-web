package main

import (
	"encoding/json"
	"github.com/JoseObreque/go-web/cmd/server/handler"
	"github.com/JoseObreque/go-web/internal/domain"
	"github.com/JoseObreque/go-web/internal/product"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	// Load environment variables
	err := godotenv.Load("./cmd/local.env")
	if err != nil {
		panic(err)
	}

	// Extract products data from the JSON file
	var productList []domain.Product
	err = extractData("products.json", &productList)
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

/*
This function extracts the data from the given JSON file and stores it in the given slice.
It returns an error if the file cannot be opened.
*/
func extractData(filename string, target *[]domain.Product) error {
	// read the JSON file
	data, err := os.ReadFile(filename)

	// If there is an error opening the file, is returned
	if err != nil {
		return err
	}

	// Unmarshal the JSON into the global products slice
	if err := json.Unmarshal(data, target); err != nil {
		return err
	}

	return nil
}
