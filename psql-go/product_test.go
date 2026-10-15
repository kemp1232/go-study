package main

import (
	"context"
	"fmt"
	"testing"
)

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
			}

			if err := app.setupDBConnection(true); err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() {
					app.pool.Close()
				})


            productRequest := ProductRequest{
                Name:     "CBR",
                Price:    test.price,
                Quantity: test.stock,
            }


            product, err := app.addProduct(context.Background(), productRequest)

            if test.wantErr && err == nil {
                t.Error("Expected validation error")
            }

			if !test.wantErr && err != nil {
				errMsg := fmt.Sprintf("Did not expected an error but got: %v", err)
				t.Error(errMsg)
			}

			if !test.wantErr && err == nil {
				if product.ID <= 0 {
					t.Error("Expected a generated product ID")
				}

				if product.Name != productRequest.Name {
					t.Errorf("Expected name %q, got %q",
						productRequest.Name, product.Name)
				}

				if product.Price != productRequest.Price {
					t.Errorf("Expected price %v, got %v",
						productRequest.Price, product.Price)
				}

				if product.Quantity != productRequest.Quantity {
					t.Errorf("Expected quantity %d, got %d",
						productRequest.Quantity, product.Quantity)
				}
			}

			if product.ID > 0 {
				t.Cleanup(func() {
					_, err := app.pool.Exec(
						context.Background(),
						"DELETE FROM products WHERE id = $1",
						product.ID,
					)
					if err != nil {
						t.Errorf("cleanup product: %v", err)
					}
				})

				var savedProduct Product
				err = app.pool.QueryRow(
					context.Background(),
					`SELECT id, name, category, price, quantity
					FROM products WHERE id = $1`,
					product.ID,
				).Scan(
					&savedProduct.ID,
					&savedProduct.Name,
					&savedProduct.Category,
					&savedProduct.Price,
					&savedProduct.Quantity,
				)

				if err != nil {
					t.Fatal(err)
				}

				if savedProduct != product {
					t.Errorf("Expected %+v, got %+v", product, savedProduct)
				}
			}
		})
	}
}

func TestListProduct(t *testing.T) {
	app := App{
	}

	if err := app.setupDBConnection(true); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		app.pool.Close()
	})


	productRequest := ProductRequest{
		Name:     "CBR",
		Price:    1000,
		Quantity: 10,
	}


	product, err := app.addProduct(context.Background(), productRequest)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	productList, err := app.listProducts(context.Background())

	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}

	found := false

	for _, p := range productList {
		if p.ID == product.ID {
			found = true

			if p != product {
				t.Errorf("Expected %+v, got %+v", product, p)
			}

			break
		}
	}

	if !found {
		t.Errorf("Product ID %d not found in list", product.ID)
	}

	if product.ID > 0 {
		t.Cleanup(func() {
			_, err := app.pool.Exec(
				context.Background(),
				"DELETE FROM products WHERE id = $1",
				product.ID,
			)
			if err != nil {
				t.Errorf("cleanup product: %v", err)
			}
		})

		var savedProduct Product
		err = app.pool.QueryRow(
			context.Background(),
			`SELECT id, name, category, price, quantity
			FROM products WHERE id = $1`,
			product.ID,
		).Scan(
			&savedProduct.ID,
			&savedProduct.Name,
			&savedProduct.Category,
			&savedProduct.Price,
			&savedProduct.Quantity,
		)

		if err != nil {
			t.Fatal(err)
		}

		if savedProduct != product {
			t.Errorf("Expected %+v, got %+v", product, savedProduct)
		}
	}
}