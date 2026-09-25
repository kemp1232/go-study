package main

import (
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
		ID: 123,
		Name: "CBR 650R",
		Price: 8000,
		Stock: 20,
	}

	product, err := app.getProductByID(123)
	if err != nil {
		t.Errorf("Did not Expect an error, got %v", err)
	}

	if product == (Product{}){
		t.Errorf("Expected product to have value, got %v", product)
	}
}