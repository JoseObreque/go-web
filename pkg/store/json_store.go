package store

import (
	"encoding/json"
	"errors"
	"github.com/JoseObreque/go-web/internal/domain"
	"os"
)

/*
The Store interface defines methods for interact with a JSON file of Products.
*/
type Store interface {
	Load() ([]domain.Product, error)
	Save([]domain.Product) error
	GetAll() ([]domain.Product, error)
	GetOne(id int) (domain.Product, error)
	AddOne(product domain.Product) error
	UpdateOne(updatedProduct domain.Product) error
	DeleteOne(id int) error
}

// The jsonStore struct is the implementation of the Store interface.
type jsonStore struct {
	filepath string
}

// NewJsonStore is a constructor for a new jsonStore instance.
func NewJsonStore(filepath string) Store {
	return &jsonStore{
		filepath: filepath,
	}
}

// The Load method retrieves all the products from a JSON file as a slice of Products.
func (s *jsonStore) Load() ([]domain.Product, error) {
	// Read all the data from the JSON file
	var products []domain.Product
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return products, err
	}

	// Unmarshal the data into a slice of Product structs
	if err = json.Unmarshal(data, &products); err != nil {
		return products, err
	}

	return products, nil
}

// The Save method saves all the products in a JSON file.
func (s *jsonStore) Save(products []domain.Product) error {
	// Marshal the data into a JSON format
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}

	// Write the data to the JSON file
	return os.WriteFile(s.filepath, data, 0644)
}

// The GetAll method retrieves all the products from a JSON file as a slice of Products.
func (s *jsonStore) GetAll() ([]domain.Product, error) {
	// Read all the data from a JSON file using the Load method
	products, err := s.Load()
	if err != nil {
		return products, err
	}

	return products, nil
}

// The GetOne method retrieves a single product from a JSON file.
func (s *jsonStore) GetOne(id int) (domain.Product, error) {
	// Read all the data from a JSON file using the Load method
	products, err := s.Load()
	if err != nil {
		return domain.Product{}, err
	}

	// Search for a product matching the ID specified
	for _, product := range products {
		if product.Id == id {
			return product, nil
		}
	}

	// If no product was found, return an error
	return domain.Product{}, errors.New("product not found")
}

// The AddOne method adds a single product to a JSON file.
func (s *jsonStore) AddOne(product domain.Product) error {
	// Load the data from a JSON file using the Load method
	products, err := s.Load()
	if err != nil {
		return err
	}

	// Update the product id and append it in the slice
	product.Id = len(products) + 1
	products = append(products, product)

	// Save the data to the JSON file
	return s.Save(products)
}

// The UpdateOne method updates a single product in a JSON file.
func (s *jsonStore) UpdateOne(updatedProduct domain.Product) error {
	// Load the data from a JSON file using the Load method
	products, err := s.Load()
	if err != nil {
		return err
	}

	// Search for a product matching the ID specified
	for i, product := range products {
		if product.Id == updatedProduct.Id {
			products[i] = updatedProduct
			return s.Save(products)
		}
	}

	// If no product was found, return an error
	return errors.New("product not found")
}

// The DeleteOne method deletes a single product from a JSON file.
func (s *jsonStore) DeleteOne(id int) error {
	// Load the data from a JSON file using the Load method
	products, err := s.Load()
	if err != nil {
		return err
	}

	// Search for a product matching the ID specified
	for i, product := range products {
		if product.Id == id {
			products = append(products[:i], products[i+1:]...)
			return s.Save(products)
		}
	}

	// If no product was found, return an error
	return errors.New("product not found")
}
