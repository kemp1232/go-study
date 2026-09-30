package main

import (
	"fmt"
	"testing"
)

func TestGetProductbyIDNotExist(t *testing.T) {
	app := App{}

	product, err := app.getProductByID(1)
	if err == nil {
		t.Error("Expected an error for non-existing order id, got nil")
	}

	if product != (Product{}) {
		t.Errorf("Expected product to be empty, got %v", product)
	}
}

func TestGetProductbyID(t *testing.T) {
	app := App{
		products: map[int]Product{},
	}

	app.products[123] = Product{
		ID:    123,
		Name:  "CBR 650R",
		Price: 8000,
		Stock: 20,
	}

	product, err := app.getProductByID(123)
	if err != nil {
		t.Errorf("Did not Expect an error, got %v", err)
	}

	if product == (Product{}) {
		t.Errorf("Expected product to have value, got %v", product)
	}
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name    string
		price   float64
		stock   int
		wantErr bool
	}{
		{"Invalid Price", 0, 10, true},
		{"Invalid stock", 23, -1, true},
		{"Add Product", 23, 11, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := App{
				products: map[int]Product{},
			}

			product, err := app.addProduct(
				"CBR650R",
				test.price,
				test.stock,
			)

			if test.wantErr && err == nil {
				t.Error("Expected an error but got nil")
			}

			if !test.wantErr && err != nil {
				errMsg := fmt.Sprintf("Did not expected an error but got: %v", err)
				t.Error(errMsg)
			}

			if !test.wantErr && err == nil {
				if product == (Product{}) {
					t.Error("Successfully added product but did not return the created product")
				}
			}
		})
	}
}
