package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var products []Product // Slice that will contain all the available products
var autoId int         // Auto-incremental id for products

func main() {
	// Get the data from the JSON file
	err := extractData("products.json", &products)
	if err != nil {
		panic(err)
	}

	// Assign the last id value to the autoId variable
	autoId = products[len(products)-1].Id

	// Create a gin router
	router := gin.Default()

	// Ungrouped endpoints
	router.GET("/ping", pongHandler)

	// Product endpoints
	productGroup := router.Group("/products")
	{
		productGroup.GET("/all", getAllProductsHandler)
		productGroup.GET("/:id", getProductByIdHandler)       // using parameters
		productGroup.GET("/search", getProductByPriceHandler) //using query parameters
		productGroup.POST("/new", createProductHandler)
	}

	// Start the server
	err = router.Run(":8080")
	if err != nil {
		return
	}
}

/*
----------------------------------------------------------------
PRODUCT DEFINITION
----------------------------------------------------------------
*/

// Product represents a single product from the supermarket.
type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	CodeValue   string  `json:"code_value" binding:"required"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

/*
------------------------------------------------------------------------------
AUXILIARY FUNCTIONS
------------------------------------------------------------------------------
*/

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

// This function checks if the given date string is a valid date.
func isValidDate(date string) bool {
	parsedDate, err := time.Parse("02/01/2006", date)
	if err != nil {
		return false
	}

	if err == nil && parsedDate.After(time.Now()) {
		return true
	}

	return false
}

/*
------------------------------------------------------------------------------
HANDLERS
------------------------------------------------------------------------------
*/

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
	var ()
	// Obtains the product id from a query
	idString := c.Param("id")

	// Parse the id string into an int
	id, err := strconv.Atoi(idString)

	// If the id is not valid return a 400 status code
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
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

// Handles the request for creating a new product.
func createProductHandler(c *gin.Context) {
	var newProduct Product

	// Obtains the new product data from the request body
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Checks if the product expiration date is valid (DD/MM/YYYY)
	if !isValidDate(newProduct.Expiration) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product expiration date is invalid"})
		return
	}

	// Checks if the code value is unique
	for _, storedProduct := range products {
		if newProduct.CodeValue == storedProduct.CodeValue {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product code cannot be duplicated"})
			return
		}
	}

	// Increment the autoId variable value and assign it to the new product
	autoId++
	newProduct.Id = autoId

	// Store the new product in the products slice
	products = append(products, newProduct)

	// Return a 200 status code with the new product
	c.JSON(http.StatusCreated, newProduct)
}
