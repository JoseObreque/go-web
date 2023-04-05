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

var (
	ErrInvalidId    = errors.New("invalid product id")
	ErrInvalidPrice = errors.New("invalid product price")
	ErrInvalidData  = errors.New("invalid product data")
)

// ProductHandler is a handler for the product endpoints.
type ProductHandler struct {
	service product.Service
}

/*
The NewProductHandler function returns a new ProductHandler. It uses the provided service for
make CRUD operations for products.
*/
func NewProductHandler(service product.Service) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

/*
The GetAll method returns all available products. It returns a HandlerFunc that
can be used to handle a GET request from the client for retrieving all products.
*/
func (h *ProductHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products := h.service.GetAll()
		c.JSON(http.StatusOK, products)
	}
}

/*
The GetById method returns a HandlerFunc that can be used to handle a GET request
from the client for retrieving a single product based on its ID (sent as URL parameter).
*/
func (h *ProductHandler) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidId.Error()})
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
The GetByPriceGt method returns a HandlerFunc that can be used to handle a GET request
from the client for retrieving all products that have a price greater than the provided
(as query parameter).
*/
func (h *ProductHandler) GetByPriceGt() gin.HandlerFunc {
	return func(c *gin.Context) {
		stringPriceGt := c.Query("priceGt")
		priceGt, err := strconv.ParseFloat(stringPriceGt, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidPrice.Error()})
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
The Create method is used to create a new product. It returns a HandlerFunc that
can be used to handle a POST request from the client for product creation.
*/
func (h *ProductHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newProduct domain.Product

		// Obtains the new product data from the request body
		if err := c.ShouldBindJSON(&newProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidData.Error()})
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
The FullUpdate method is used to update a product. It returns a HandlerFunc that
can be used to handle a PUT request from the client for updating a product.
*/
func (h *ProductHandler) FullUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtains the product id from a URL parameter
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidId.Error()})
			return
		}

		// Extract the product data from the request body
		var newProductData domain.Product
		if err := c.ShouldBindJSON(&newProductData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidData.Error()})
			return
		}
		// Checks if the product expiration date is valid (DD/MM/YYYY)
		isValidDate, err := validateDate(newProductData.Expiration)
		if !isValidDate {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Updates the product
		updatedProduct, err := h.service.Update(id, newProductData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedProduct)
	}
}

/*
The PartialUpdate method is used to update some fields of a product. It returns a HandlerFunc
that can be used to handle a PUT request from the client for partially updating a product.
*/
func (h *ProductHandler) PartialUpdate() gin.HandlerFunc {
	type Request struct {
		Name        string  `json:"name,omitempty"`
		Quantity    int     `json:"quantity,omitempty"`
		CodeValue   string  `json:"code_value,omitempty"`
		IsPublished bool    `json:"is_published,omitempty"`
		Expiration  string  `json:"expiration,omitempty"`
		Price       float64 `json:"price,omitempty"`
	}
	return func(c *gin.Context) {
		// Obtains the product id from a URL parameter
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidId.Error()})
			return
		}

		// Extract the product data from the request body
		var partialUpdateData Request
		if err := c.ShouldBindJSON(&partialUpdateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidData.Error()})
			return
		}

		update := domain.Product{
			Name:        partialUpdateData.Name,
			Quantity:    partialUpdateData.Quantity,
			CodeValue:   partialUpdateData.CodeValue,
			IsPublished: partialUpdateData.IsPublished,
			Expiration:  partialUpdateData.Expiration,
			Price:       partialUpdateData.Price,
		}

		// Checks if the product expiration date is valid (DD/MM/YYYY)
		if update.Expiration != "" {
			isValidDate, err := validateDate(update.Expiration)
			if !isValidDate {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		// Updates the product
		updatedProduct, err := h.service.Update(id, update)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedProduct)
	}
}

/*
The Delete method is used to delete a product. It returns a HandlerFunc that
can be used to handle a DELETE request from the client for deleting a product.
*/
func (h *ProductHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtains the product id from a URL parameter
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidId.Error()})
			return
		}

		// Deletes the product
		err = h.service.Delete(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
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
