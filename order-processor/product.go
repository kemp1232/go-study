package main

import (
	"errors"
	"fmt"
)

type Product struct {
	ID int	
	Name string	
	Price float64	
	Stock int	
}

func (app *App) addProduct(name string, price float64, stock int) {
	id := app.nextProductID

	newProduct := Product{
		ID: id,
		Name: name,
		Price: price,
		Stock: stock,
	}

	app.products[id] = newProduct
	app.nextProductID++
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