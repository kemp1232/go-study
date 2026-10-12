package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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

			app.setupDBConnection()
			defer app.conn.Close(context.Background())

            productRequest := ProductRequest{
                Name:     "CBR",
                Price:    test.price,
                Quantity: test.stock,
            }

			req := httptest.NewRequest(
                http.MethodPost,
                "/products",
                nil,
            )

            product, err := app.addProduct(req.Context(), productRequest)

            if test.wantErr && err == nil {
                t.Error("Expected validation error")
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