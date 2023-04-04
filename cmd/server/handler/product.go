package handler

import (
	"errors"
	"github.com/JoseObreque/go-web/internal/domain"
	"github.com/JoseObreque/go-web/internal/product"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// ProductHandler is a handler for the product endpoints.
type ProductHandler struct {
	service product.Service
}

/*
NewProductHandler returns a new ProductHandler. It uses the provided service for
make CRUD operations for products.
*/
func NewProductHandler(service product.Service) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

/*
The GetAll function returns all available products. It returns a HandlerFunc that
can be used to handle a GET request from the client for retrieving all products.
*/
func (h *ProductHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products := h.service.GetAll()
		c.JSON(http.StatusOK, products)
	}
}

/*
The GetById function returns a HandlerFunc that can be used to handle a GET request
from the client for retrieving a single product based on its ID (sent as URL parameter).
*/
func (h *ProductHandler) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
			return
		}

		targetProduct, err := h.service.GetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, targetProduct)
	}
}

/*
The GetByPriceGt function returns a HandlerFunc that can be used to handle a GET request
from the client for retrieving all products that have a price greater than the provided
(as query parameter).
*/
func (h *ProductHandler) GetByPriceGt() gin.HandlerFunc {
	return func(c *gin.Context) {
		stringPriceGt := c.Query("priceGt")
		priceGt, err := strconv.ParseFloat(stringPriceGt, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid product price")})
			return
		}

		filteredProducts, err := h.service.GetByPriceGt(priceGt)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, filteredProducts)
	}
}

/*
The Create function is used to create a new product. It returns a HandlerFunc that
can be used to handle a POST request from the client for product creation.
*/
func (h *ProductHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newProduct domain.Product

		// Obtains the new product data from the request body
		if err := c.ShouldBindJSON(&newProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Checks if the product expiration date is valid (DD/MM/YYYY)
		validDate, err := validateDate(newProduct.Expiration)
		if !validDate {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Creates the new product
		createdProduct, err := h.service.Create(newProduct)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdProduct)
	}
}

/*
A function that checks if a given date string is a valid date. It returns true if the
date string is a valid date and occurs after the current date. Otherwise, it returns false with
an error.
*/
func validateDate(date string) (bool, error) {
	parsedDate, err := time.Parse("02/01/2006", date)
	if err != nil {
		return false, errors.New("invalid expiration date format")
	}

	if err == nil && parsedDate.Before(time.Now()) {
		return false, errors.New("expiration date must be after current date")
	}

	return true, nil
}
