package handler

import (
	"errors"
	"github.com/JoseObreque/go-web/internal/domain"
	"github.com/JoseObreque/go-web/internal/product"
	"github.com/JoseObreque/go-web/pkg/web"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	ErrInvalidId    = errors.New("invalid product id")
	ErrInvalidPrice = errors.New("invalid product price")
	ErrInvalidData  = errors.New("invalid product data")
	ErrNotFound     = errors.New("product not found")
	ErrInvalidCode  = errors.New("invalid product code value")
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

// GetAll godoc
// @Summary List all products
// @Tags Products
// @Description List all available products
// @Produce json
// @Success 200 {object} web.Response
// @Router /products/all [get]
func (h *ProductHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products := h.service.GetAll()
		web.Success(c, 200, products)
	}
}

// GetById godoc
// @Summary Get a specific product
// @Tags Products
// @Description Get a specific product based on its ID
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			web.Failure(c, 400, ErrInvalidId)
			return
		}

		targetProduct, err := h.service.GetById(id)
		if err != nil {
			web.Failure(c, 404, err)
			return
		}

		web.Success(c, 200, targetProduct)
	}
}

// GetByPriceGt godoc
// @Summary Get all products based on its price
// @Tags Products
// @Description Get all products with a price greater than the provided value
// @Produce json
// @Param priceGt query int true "Price"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Router /products/search [get]
func (h *ProductHandler) GetByPriceGt() gin.HandlerFunc {
	return func(c *gin.Context) {
		stringPriceGt := c.Query("priceGt")
		priceGt, err := strconv.ParseFloat(stringPriceGt, 64)
		if err != nil {
			web.Failure(c, 400, ErrInvalidPrice)
			return
		}

		filteredProducts, err := h.service.GetByPriceGt(priceGt)
		if err != nil {
			web.Failure(c, 404, err)
			return
		}

		web.Success(c, 200, filteredProducts)
	}
}

// Create godoc
// @Summary Create a new product
// @Tags Products
// @Description Create a new product and store it in the database
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param newProduct body domain.ProductRequest true "new product"
// @Success 201 {object} web.Response
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Router /products/new [post]
func (h *ProductHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtains the new product data from the request body
		var newProduct domain.Product
		if err := c.ShouldBindJSON(&newProduct); err != nil {
			web.Failure(c, 400, ErrInvalidData)
			return
		}

		// Checks if the product expiration date is valid (DD/MM/YYYY)
		validDate, err := validateDate(newProduct.Expiration)
		if !validDate {
			web.Failure(c, 400, err)
			return
		}

		// Creates the new product
		createdProduct, err := h.service.Create(newProduct)
		if err != nil {
			web.Failure(c, 400, err)
			return
		}

		web.Success(c, 201, createdProduct)
	}
}

// FullUpdate godoc
// @Summary Update a product
// @Tags Products
// @Description Update all the fields of a product.
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param id path int true "Product ID"
// @Param partialUpdateData body domain.ProductRequest true "updated product"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) FullUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Checks if the given token is valid
		err := isAuthorized(c)
		if err != nil {
			web.Failure(c, 401, err)
			return
		}

		// Obtains the product id from a URL parameter
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			web.Failure(c, 400, ErrInvalidId)
			return
		}

		// Extract the product data from the request body
		var newProductData domain.Product
		if err := c.ShouldBindJSON(&newProductData); err != nil {
			web.Failure(c, 400, ErrInvalidData)
			return
		}
		// Checks if the product expiration date is valid (DD/MM/YYYY)
		isValidDate, err := validateDate(newProductData.Expiration)
		if !isValidDate {
			web.Failure(c, 400, err)
			return
		}

		// Updates the product
		updatedProduct, err := h.service.Update(id, newProductData)

		// Check for errors
		if err != nil && err.Error() == ErrNotFound.Error() {
			web.Failure(c, 404, err)
			return
		}

		if err != nil && err.Error() == ErrInvalidCode.Error() {
			web.Failure(c, 400, err)
			return
		}

		web.Success(c, 200, updatedProduct)
	}
}

// PartialUpdate godoc
// @Summary Partially update a product
// @Tags Products
// @Description Update some product fields data
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param id path int true "Product ID"
// @Param partialUpdateData body domain.ProductRequest true "updated product"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Router /products/{id} [patch]
func (h *ProductHandler) PartialUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Checks if the given token is valid
		err := isAuthorized(c)
		if err != nil {
			web.Failure(c, 401, err)
			return
		}

		// Obtains the product id from a URL parameter
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			web.Failure(c, 400, ErrInvalidId)
			return
		}

		// Extract the product data from the request body
		var partialUpdateData domain.ProductRequest
		if err := c.ShouldBindJSON(&partialUpdateData); err != nil {
			web.Failure(c, 400, ErrInvalidData)
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
				web.Failure(c, 400, err)
				return
			}
		}

		// Updates the product
		updatedProduct, err := h.service.Update(id, update)

		// Check for errors
		if err != nil && err.Error() == ErrNotFound.Error() {
			web.Failure(c, 404, err)
			return
		}
		if err != nil && err.Error() == ErrInvalidCode.Error() {
			web.Failure(c, 400, err)
			return
		}

		web.Success(c, 200, updatedProduct)
	}
}

// Delete godoc
// @Summary Delete a product
// @Tags Products
// @Description Delete permanently a product
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param id path int true "Product ID"
// @Success 204 {object} web.Response
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Checks if the given token is valid
		err := isAuthorized(c)
		if err != nil {
			web.Failure(c, 401, err)
			return
		}

		// Obtains the product id from a URL parameter
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			web.Failure(c, 400, ErrInvalidId)
			return
		}

		// Deletes the product
		err = h.service.Delete(id)
		if err != nil {
			web.Failure(c, 404, err)
			return
		}

		web.Success(c, http.StatusNoContent, nil)
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

// Auxiliary function that checks if the given token is valid.
func isAuthorized(c *gin.Context) error {
	// Get the token from the header
	token := c.GetHeader("token")

	// Authentication
	if token != os.Getenv("TOKEN") {
		return errors.New("invalid token")
	}
	return nil
}
