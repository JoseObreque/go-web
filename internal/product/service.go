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

// GetAll returns all available products
func (s *ServiceImpl) GetAll() []domain.Product {
	return s.repository.GetAll()
}

// GetById returns a product by its ID
func (s *ServiceImpl) GetById(id int) (domain.Product, error) {
	product, err := s.repository.GetById(id)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

/*
GetByPriceGt returns all product that has a price greater than the given price.
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
The Create function try to create a new product. If the product already exists, it returns an error.
Otherwise, it creates a new product and returns it.
*/
func (s *ServiceImpl) Create(product domain.Product) (domain.Product, error) {
	newProduct, err := s.repository.Create(product)
	if err != nil {
		return domain.Product{}, err
	}
	return newProduct, nil
}
