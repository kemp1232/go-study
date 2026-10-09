package main

import (
	"errors"
)
type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}



func (app *App) listProducts() map[int]Product {
	return app.products
}

func (app *App) getProductByID(productID int) (Product, error) {
	product, exists := app.products[productID]
	if !exists {
		return Product{}, ErrNotFound
	}

	return product, nil
}

func (app *App) addProduct(productRequest ProductRequest) (Product, error) {
	if productRequest.Price <= 0 {
		return Product{}, errors.New("Price must be greater than 0.")
	}

	if productRequest.Quantity < 0 {
		return Product{}, errors.New("Quantity must not be negative.")
	}

	id := app.nextProductID
	newProduct := Product{
		ID:    id,
		Name:  productRequest.Name,
		Category:  productRequest.Category,
		Price: productRequest.Price,
		Quantity: productRequest.Quantity,
	}

	app.products[id] = newProduct
	app.nextProductID++
	return newProduct, nil
}

func (app *App) editProduct(productID int, productRequest ProductRequest) (Product, error) {
	product, exists := app.products[productID]

	if !exists {
		return Product{}, ErrNotFound
	}
	
	if productRequest.Price <= 0 {
		return Product{}, errors.New("Price must be greater than 0.")
	}

	if productRequest.Quantity < 0 {
		return Product{}, errors.New("Quantity must not be negative.")
	}

	updatedProduct := Product{
		ID:    product.ID,
		Name:  productRequest.Name,
		Category:  productRequest.Category,
		Price: productRequest.Price,
		Quantity: productRequest.Quantity,
	}

	app.products[product.ID] = updatedProduct
	return updatedProduct, nil
}


func (app *App) deleteProduct(productID int) (bool, error) {
	product, exists := app.products[productID]

	if !exists {
		return false, ErrNotFound
	}
	
	delete(app.products, product.ID)
	return true, nil
}