package main

import (
	"encoding/json"
	"os"
)

// Slice that will contain all the available products
var products []Product

func main() {

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
