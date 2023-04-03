package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Slice that will contain all the available products
var products []Product

func main() {
	// Get the data from the JSON file
	err := extractData("products.json", &products)
	if err != nil {
		panic(err)
	}

	// Create a gin router
	router := gin.Default()

	// Endpoints
	router.GET("/ping", pongHandler)

	// Product endpoints
	productGroup := router.Group("/products")
	productGroup.GET("/all", getAllProductsHandler)
	productGroup.GET("/:id", getProductByIdHandler)
	productGroup.GET("/search", getProductByPriceHandler)

	// Start the server
	err = router.Run(":8080")
	if err != nil {
		return
	}
}

// Product represents a product from the supermarket.
type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

/*
This function extracts the data from the given JSON file and stores it in the given slice.
It returns an error if the file cannot be opened.
*/
func extractData(filename string, target *[]Product) error {
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

// Handles the request for the /ping endpoint. It returns a simple text response ("pong").
func pongHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// Handles the request for all the available products.
func getAllProductsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

// Handles the request for a single product based on its ID.
func getProductByIdHandler(c *gin.Context) {
	// Obtains the product id from a query
	idString := c.Param("id")

	// Parse the id string into an int
	id, err := strconv.Atoi(idString)

	// If the id is not valid return a 400 status code
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Iterate over the products slice and find the product with the given id
	for _, product := range products {
		if product.Id == id {
			c.JSON(http.StatusOK, product)
			return
		}
	}

	// Return a 404 status code if the product is not found
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

// Handles the request for all products that have a price greater than the given price.
func getProductByPriceHandler(c *gin.Context) {
	// Initialize an empty slice of filtered products
	var filteredProducts []Product

	// Obtains the price from a query
	priceGtString := c.Query("priceGt")

	// Parse the price string into a float64
	priceGt, err := strconv.ParseFloat(priceGtString, 64)

	// If the price is not valid return a 400 status code
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error de parseo"})
		fmt.Println("error parsing")
		return
	}

	// Iterate over the products slice and find the product with the given price
	for _, product := range products {
		if product.Price > priceGt {
			filteredProducts = append(filteredProducts, product)
		}
	}

	// Return a 404 status code if no products were found,
	// otherwise return a 200 status code with the products filtered.
	if len(filteredProducts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	} else {
		c.JSON(http.StatusOK, filteredProducts)
	}
}
