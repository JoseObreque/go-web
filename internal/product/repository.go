package product

import (
	"errors"
	"github.com/JoseObreque/go-web/internal/domain"
)

// Repository is the interface definition for the product service
type Repository interface {
	GetAll() []domain.Product
	GetById(id int) (domain.Product, error)
	GetByPriceGt(price float64) []domain.Product
	Create(product domain.Product) (domain.Product, error)
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

// GetAll returns all available products
func (r *RepositoryImpl) GetAll() []domain.Product {
	return r.productList
}

// GetById returns a product by its ID
func (r *RepositoryImpl) GetById(id int) (domain.Product, error) {
	for _, product := range r.productList {
		if product.Id == id {
			return product, nil
		}
	}

	return domain.Product{}, errors.New("product not found")
}

// GetByPriceGt returns a list of products with a price greater than the given price
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
Create function creates a new product. If the product code already exists, it will return an error.
Otherwise, it creates a new product.
*/
func (r *RepositoryImpl) Create(product domain.Product) (domain.Product, error) {
	if !r.validateCodeValue(product.CodeValue) {
		return domain.Product{}, errors.New("invalid code value")
	}

	product.Id = len(r.productList) + 1
	r.productList = append(r.productList, product)

	return product, nil
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
