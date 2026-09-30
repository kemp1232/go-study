package main

import (
	"errors"
	"fmt"
)

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}

func (app *App) addProduct(name string, price float64, stock int) (Product, error) {
	if price <= 0 {
		return Product{}, errors.New("Price must be greater than 0.")
	}

	if stock < 0 {
		return Product{}, errors.New("Stock must not be negative.")
	}

	id := app.nextProductID
	newProduct := Product{
		ID:    id,
		Name:  name,
		Price: price,
		Stock: stock,
	}

	app.products[id] = newProduct
	app.nextProductID++
	return newProduct, nil
}

func (app *App) listProducts() {
	for _, product := range app.products {
		fmt.Printf(
			"%v %v $%v %v Pcs \n",
			product.ID,
			product.Name,
			product.Price,
			product.Stock,
		)
	}
}

func (app *App) getProductByID(productID int) (Product, error) {
	product, exists := app.products[productID]
	if !exists {
		errMsg := fmt.Sprintf("Product ID: %v does not exist:", productID)
		return Product{}, errors.New(errMsg)
	}

	return product, nil
}
