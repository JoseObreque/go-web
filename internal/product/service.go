package product

import (
	"errors"
	"github.com/JoseObreque/go-web/internal/domain"
)

type Service interface {
	GetAll() []domain.Product
	GetById(id int) (domain.Product, error)
	GetByPriceGt(price float64) ([]domain.Product, error)
	Create(product domain.Product) (domain.Product, error)
	Update(id int, updatedProduct domain.Product) (domain.Product, error)
	Delete(id int) error
}

type ServiceImpl struct {
	repository Repository
}

// The NewService function returns a new instance of the service.
func NewService(repository Repository) Service {
	return &ServiceImpl{
		repository: repository,
	}
}

// The GetAll method returns all available products
func (s *ServiceImpl) GetAll() []domain.Product {
	return s.repository.GetAll()
}

// The GetById method returns a product by its ID
func (s *ServiceImpl) GetById(id int) (domain.Product, error) {
	product, err := s.repository.GetById(id)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

/*
The GetByPriceGt method returns all product that has a price greater than the given price.
If no product has a price greater than the given price, it returns an error.
Otherwise, it returns all product that has a price greater than the given price.
*/
func (s *ServiceImpl) GetByPriceGt(price float64) ([]domain.Product, error) {
	products := s.repository.GetByPriceGt(price)
	if len(products) == 0 {
		return []domain.Product{}, errors.New("no products found")
	}
	return products, nil
}

/*
The Create method try to create a new product. If the product already exists, it returns an error.
Otherwise, it creates a new product and returns it.
*/
func (s *ServiceImpl) Create(product domain.Product) (domain.Product, error) {
	newProduct, err := s.repository.Create(product)
	if err != nil {
		return domain.Product{}, err
	}
	return newProduct, nil
}

/*
The Update method try to update a product. If the product does not exist or any updated fields
data is invalid then returns an error. Otherwise, it updates the product and returns it.
*/
func (s *ServiceImpl) Update(id int, newProductData domain.Product) (domain.Product, error) {
	// Search the old product data
	product, err := s.repository.GetById(id)
	if err != nil {
		return domain.Product{}, err
	}

	// Update the product data
	if newProductData.Name != "" {
		product.Name = newProductData.Name
	}
	if newProductData.Quantity > 0 {
		product.Quantity = newProductData.Quantity
	}
	if newProductData.CodeValue != "" {
		product.CodeValue = newProductData.CodeValue
	}
	if newProductData.Expiration != "" {
		product.Expiration = newProductData.Expiration
	}
	if newProductData.Price > 0 {
		product.Price = newProductData.Price
	}
	product.IsPublished = newProductData.IsPublished

	// Store the updated product data
	updatedProduct, err := s.repository.Update(id, product)
	if err != nil {
		return domain.Product{}, err
	}
	return updatedProduct, nil
}

/*
The Delete method try to delete a product. If the product does not exist, it returns an error.
*/
func (s *ServiceImpl) Delete(id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
