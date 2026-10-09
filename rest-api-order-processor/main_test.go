package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleGetProduct(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/products/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	var product Product

	err := json.NewDecoder(response.Body).Decode(&product)

	if err != nil {
		t.Fatal(err)
	}

	if product.ID != 1 {
		t.Errorf("expected ID 1, got %d", product.ID)
	}

	if product.Name != "Keyboard" {
		t.Errorf("expected Keyboard, got %s", product.Name)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}
}

func TestHandleGetProductNotFound(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/products/2",
		nil,
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()


	if response.StatusCode != http.StatusNotFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			response.StatusCode,
		)
	}
}

func TestHandleGetProductInvalidID(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/products/abc",
		nil,
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()


	if response.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			response.StatusCode,
		)
	}
}

func TestHandleGetProductList(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/products",
		nil,
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	fmt.Println("response", response)
	if response.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}
}

func TestHandleCreateProduct(t *testing.T) {
	app := App{
		products:      map[int]Product{},
		orders:        map[int]Order{},
		nextProductID: 1,
		nextOrderID:   1,
	}

	body := `{
		"name": "Keyboard",
		"category": "Electronics",
		"price": 2500,
		"quantity": 5
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)
	
	response := recorder.Result()
	
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusCreated,
			response.StatusCode,
		)
	}

	product := app.products[1]
	if product.Name != "Keyboard" {
		t.Errorf(
			"expected product name Keyboard, got %s",
			product.Name,
		)
	}

	if product.Category != "Electronics" {
		t.Errorf(
			"expected product category Electronics, got %s",
			product.Category,
		)
	}
}

func TestHandleCreateProductInvalidBody(t *testing.T) {
	app := App{
		products:      map[int]Product{},
		orders:        map[int]Order{},
		nextProductID: 1,
		nextOrderID:   1,
	}

	body := `{"name": "Keyboard", "price": }`

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)
	
	response := recorder.Result()
	
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.StatusCode,
		)
	}
}

func TestHandleCreateProductNegativePrice(t *testing.T) {
	app := App{
		products:      map[int]Product{},
		orders:        map[int]Order{},
		nextProductID: 1,
		nextOrderID:   1,
	}

	body := `{
		"name": "Keyboard",
		"category": "Electronics",
		"price": -2500,
		"quantity": 5
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)
	
	response := recorder.Result()
	
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.StatusCode,
		)
	}
}

func TestHandleUpdateProduct(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
		orders:        map[int]Order{},
		nextProductID: 1,
		nextOrderID:   1,
	}

	body := `{
		"name": "Keyboard",
		"category": "Electronics",
		"price": 2500,
		"quantity": 5
	}`

	
	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	fmt.Println("response", response)
	if response.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}
}

func TestHandleDeleteProduct(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
		orders:        map[int]Order{},
		nextProductID: 1,
		nextOrderID:   1,
	}
	
	req := httptest.NewRequest(
		http.MethodDelete,
		"/products/1",
		nil,
	)


	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	fmt.Println("response", response)
	if response.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}
}

func TestHandleGetOrderList(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
		orders: map[int]Order{
			1: {
				ID:     1,
				Items:  []OrderItem{},
				Total:  123,
				Status: "pending",
			},
		},
		nextProductID: 1,
		nextOrderID:   1,
	}
	
	req := httptest.NewRequest(
		http.MethodGet,
		"/orders",
		nil,
	)
	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	fmt.Println("response", response)
	if response.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}
}

func TestHandleGetOrder(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
		orders: map[int]Order{
			1: {
				ID:     1,
				Items:  []OrderItem{},
				Total:  123,
				Status: "pending",
			},
		},
		nextProductID: 1,
		nextOrderID:   1,
	}
	
	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	fmt.Println("response", response)
	if response.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}
}

func TestHandleGetOrderInvalidID(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
		orders: map[int]Order{
			1: {
				ID:     1,
				Items:  []OrderItem{},
				Total:  123,
				Status: "pending",
			},
		},
		nextProductID: 1,
		nextOrderID:   1,
	}
	
	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/2",
		nil,
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	fmt.Println("response", response)
	if response.StatusCode != http.StatusNotFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			response.StatusCode,
		)
	}
}

func TestHandleCreateOrder(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
		orders: map[int]Order{
			1: {
				ID:     1,
				Items:  []OrderItem{},
				Total:  123,
				Status: "pending",
			},
		},
		nextProductID: 1,
		nextOrderID:   1,
	}

	body := `{
		"items": [
			{
				"productID": 1,
				"quantity": 2
			}
		]
	}`
	
	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	fmt.Println("response", response)
	if response.StatusCode != http.StatusCreated {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusCreated,
			response.StatusCode,
		)
	}
}

func TestHandleCreateOrderInvalidProductID(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
		orders: map[int]Order{
			1: {
				ID:     1,
				Items:  []OrderItem{},
				Total:  123,
				Status: "pending",
			},
		},
		nextProductID: 1,
		nextOrderID:   1,
	}

	body := `{
		"items": [
			{
				"productID": 2,
				"quantity": 2
			}
		]
	}`
	
	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	fmt.Println("response", response)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.StatusCode,
		)
	}
}

func TestHandleCompleteOrder(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
		orders: map[int]Order{
			1: {
				ID:     1,
				Items:  []OrderItem{
					{
						ProductID: 1,
						Quantity: 3,
					},
				},
				Total:  123,
				Status: "pending",
			},
		},
		nextProductID: 1,
		nextOrderID:   1,
	}


	req := httptest.NewRequest(
		http.MethodPost,
		"/orders/1/complete",
		nil,
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}

	if app.products[1].Quantity != 2 {
		t.Errorf("expected stock 2, got %d", app.products[1].Quantity)
	}

	if app.orders[1].Status != "completed" {
		t.Errorf("expected completed status, got %s", app.orders[1].Status)
	}
}

func TestHandleCompleteOrderInvalidID(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
		orders: map[int]Order{
			1: {
				ID:     1,
				Items:  []OrderItem{
					{
						ProductID: 1,
						Quantity: 3,
					},
				},
				Total:  123,
				Status: "pending",
			},
		},
		nextProductID: 1,
		nextOrderID:   1,
	}


	req := httptest.NewRequest(
		http.MethodPost,
		"/orders/2/complete",
		nil,
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()
	if response.StatusCode != http.StatusNotFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			response.StatusCode,
		)
	}
}

func TestHandleCancelOrder(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
		orders: map[int]Order{
			1: {
				ID:     1,
				Items:  []OrderItem{
					{
						ProductID: 1,
						Quantity: 3,
					},
				},
				Total:  123,
				Status: "pending",
			},
		},
		nextProductID: 1,
		nextOrderID:   1,
	}


	req := httptest.NewRequest(
		http.MethodPost,
		"/orders/1/cancel",
		nil,
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	fmt.Println("response", response)
	if response.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}

	if app.products[1].Quantity != 5 {
		t.Errorf("expected stock 5, got %d", app.products[1].Quantity)
	}

	if app.orders[1].Status != "cancelled" {
		t.Errorf("expected cancelled status, got %s", app.orders[1].Status)
	}
}

func TestHandleCancelOrderInvalidID(t *testing.T) {
	app := App{
		products: map[int]Product{
			1: {
				ID:       1,
				Name:     "Keyboard",
				Category: "Electronics",
				Price:    2500,
				Quantity: 5,
			},
		},
		orders: map[int]Order{
			1: {
				ID:     1,
				Items:  []OrderItem{
					{
						ProductID: 1,
						Quantity: 3,
					},
				},
				Total:  123,
				Status: "pending",
			},
		},
		nextProductID: 1,
		nextOrderID:   1,
	}


	req := httptest.NewRequest(
		http.MethodPost,
		"/orders/2/cancel",
		nil,
	)

	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()
	fmt.Println("response", response)
	if response.StatusCode != http.StatusNotFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			response.StatusCode,
		)
	}
}