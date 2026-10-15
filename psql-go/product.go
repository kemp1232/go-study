package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func (app *App) listProducts(r context.Context) ([]Product, error) {
	ctx, cancel := context.WithTimeout(r, 5*time.Second)
	defer cancel()
	products := []Product{}

	query := `SELECT id, name, category, price, quantity from products;`
    rows, err := app.pool.Query(ctx, query)
	if err != nil {
		fmt.Println("err", err)
		return nil, fmt.Errorf("%w: %v", ErrDBQueryIssue, err)
	}
	defer rows.Close()

	for rows.Next() {
		var product Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Category,
			&product.Price,
			&product.Quantity,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDBQueryIssue, err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBQueryIssue, err)
	}

	return products, nil
}

func (app *App) addProduct(ctx context.Context, productRequest ProductRequest) (Product, error) {
	if productRequest.Price <= 0 {
		return Product{}, errors.New("Price must be greater than 0.")
	}

	if productRequest.Quantity < 0 {
		return Product{}, errors.New("Quantity must not be negative.")
	}

	
	query := `
		INSERT INTO products (name, category, price, quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, category, price, quantity
	`

	var product Product
	err := app.pool.QueryRow(
		ctx,
		query,
		productRequest.Name,
		productRequest.Category,
		productRequest.Price,
		productRequest.Quantity,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Category,
		&product.Price,
		&product.Quantity,
	)

	if err != nil {
		return Product{}, fmt.Errorf("%w: %v", ErrDBQueryIssue, err)
	}

	return product, nil
}