package product

import (
	"errors"
	"github.com/JoseObreque/go-web/internal/domain"
)

var (
	ErrNotFound    = errors.New("product not found")
	ErrInvalidCode = errors.New("invalid product code value")
)

// Repository is the interface definition for the product service
type Repository interface {
	GetAll() []domain.Product
	GetById(id int) (domain.Product, error)
	GetByPriceGt(price float64) []domain.Product
	Create(product domain.Product) (domain.Product, error)
	Update(id int, newProductData domain.Product) (domain.Product, error)
	Delete(id int) error
}

// RepositoryImpl is the implementation of the repository interface
type RepositoryImpl struct {
	productList []domain.Product
}

// The NewRepository function returns a new instance of the repository.
func NewRepository(productList []domain.Product) Repository {
	return &RepositoryImpl{
		productList: productList,
	}
}

// The GetAll method returns all available products
func (r *RepositoryImpl) GetAll() []domain.Product {
	return r.productList
}

// The GetById method returns a product by its ID
func (r *RepositoryImpl) GetById(id int) (domain.Product, error) {
	for _, product := range r.productList {
		if product.Id == id {
			return product, nil
		}
	}

	return domain.Product{}, ErrNotFound
}

// The GetByPriceGt method returns a list of products with a price greater than the given price.
func (r *RepositoryImpl) GetByPriceGt(price float64) []domain.Product {
	var filteredProducts []domain.Product

	for _, product := range r.productList {
		if product.Price > price {
			filteredProducts = append(filteredProducts, product)
		}
	}
	return filteredProducts
}

/*
The Create method creates a new product. If the product code already exists, it will return an error.
Otherwise, it creates a new product.
*/
func (r *RepositoryImpl) Create(product domain.Product) (domain.Product, error) {
	if !r.validateCodeValue(product.CodeValue) {
		return domain.Product{}, ErrInvalidCode
	}

	product.Id = len(r.productList) + 1
	r.productList = append(r.productList, product)

	return product, nil
}

/*
The Update method updates a product. It receives the ID of the product and the updated product
data as parameters and returns the updated product if the process was successful. Otherwise, it
returns an error.
*/
func (r *RepositoryImpl) Update(id int, updatedProduct domain.Product) (domain.Product, error) {
	// Search for the product with the given ID
	for i, product := range r.productList {
		if product.Id == id {
			// Validate the updated code value
			if !r.validateCodeValue(updatedProduct.CodeValue) || product.Expiration != updatedProduct.Expiration {
				return domain.Product{}, ErrInvalidCode
			}
			// Store the updated product and return it
			updatedProduct.Id = id
			r.productList[i] = updatedProduct
			return updatedProduct, nil
		}
	}
	return domain.Product{}, ErrNotFound
}

/*
The Delete method deletes a product. It receives the ID of the product and returns an error if the
product does not exist.
*/
func (r *RepositoryImpl) Delete(id int) error {
	for i, product := range r.productList {
		if product.Id == id {
			r.productList = append(r.productList[:i], r.productList[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

/*
A function that check if a given code value already exists. If it does, the code value
is invalid and returns false. Otherwise, it returns true.
*/
func (r *RepositoryImpl) validateCodeValue(codeValue string) bool {
	for _, product := range r.productList {
		if product.CodeValue == codeValue {
			return false
		}
	}
	return true
}
