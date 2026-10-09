package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type App struct {
	products      map[int]Product
	orders        map[int]Order
	nextProductID int
	nextOrderID   int
}

type ResponseMsg struct {
	Message string `json:"message"`
}

type ProductRequest struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type OrderRequest struct {
	Items  []OrderItem `json:"items"`
}

var ErrNotFound = errors.New("resource not found")
var ErrInvalidBody = errors.New("invalid body")

func main() {
	app := App{
		products:      map[int]Product{},
		orders:        map[int]Order{},
		nextProductID: 5,
		nextOrderID:   1,
	}

	app.products[1] = Product{
		ID:    1,
		Name:  "CBR 650R",
		Category: "Moto",
		Price: 8000,
		Quantity: 20,
	}

	app.products[2] = Product{
		ID:    2,
		Name:  "Speed 400",
		Category: "Moto",
		Price: 3000,
		Quantity: 10,
	}

	app.products[3] = Product{
		ID:    3,
		Name:  "Dominar 400 UG",
		Category: "Moto",
		Price: 2000,
		Quantity: 30,
	}

	app.products[4] = Product{
		ID:    4,
		Name:  "Z650RS",
		Category: "Moto",
		Price: 3500,
		Quantity: 10,
	}


	app.handleRequest()
}

func (app *App) handleRequest() {
	mux := app.routes()
	fmt.Println("Server running at localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func (app *App) routes() http.Handler {
	mux := http.NewServeMux()

	// Product routes
	mux.HandleFunc("GET /products", app.handleListProducts)
	mux.HandleFunc("POST /products", app.handleCreateProduct)
	mux.HandleFunc("GET /products/{id}", app.handleGetProduct)
	mux.HandleFunc("PUT /products/{id}", app.handleUpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", app.handleDeleteProduct)

	// Order routes
	mux.HandleFunc("GET /orders", app.handleGetOrderList)
	mux.HandleFunc("GET /orders/{id}", app.handleGetOrder)
	mux.HandleFunc("POST /orders", app.handleCreateOrder)
	mux.HandleFunc("POST /orders/{id}/complete", app.handleOrderComplete)
	mux.HandleFunc("POST /orders/{id}/cancel", app.handleOrderCancel)
	return mux
}

func (app *App) handleListProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	products := app.listProducts()
	json.NewEncoder(w).Encode(products)
}

func (app *App) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stringID := r.PathValue("id")

	id, err := strconv.Atoi(stringID)
	if err != nil {
		ResponseMsg := ResponseMsg{
			Message: "Invalid ID.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	product, err := app.getProductByID(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			ResponseMsg := ResponseMsg{
				Message: "Product does not exist.",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseMsg)
			return
		}

		ResponseMsg := ResponseMsg{
			Message: "Something went wrong.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func (app *App) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("handleCreateProduct")

	var body ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ResponseMsg := ResponseMsg{
			Message: "Something went wrong.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	product, err := app.addProduct(body)
	if err != nil {
		ResponseMsg := ResponseMsg{
			Message: "Failed to create product",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (app *App) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("handleUpdateProduct")
	stringID := r.PathValue("id")

	id, err := strconv.Atoi(stringID)
	if err != nil {
		ResponseMsg := ResponseMsg{
			Message: "Invalid ID.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	var body ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ResponseMsg := ResponseMsg{
			Message: "Something went wrong.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	product, err := app.editProduct(id, body)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			ResponseMsg := ResponseMsg{
				Message: "Product does not exist.",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseMsg)
			return
		}

		ResponseMsg := ResponseMsg{
			Message: "Something went wrong.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func (app *App) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("handleDeleteProduct")
	stringID := r.PathValue("id")

	id, err := strconv.Atoi(stringID)
	if err != nil {
		ResponseMsg := ResponseMsg{
			Message: "Invalid ID.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	isDeleted, err := app.deleteProduct(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			ResponseMsg := ResponseMsg{
				Message: "Product does not exist.",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseMsg)
			return
		}

		ResponseMsg := ResponseMsg{
			Message: "Something went wrong.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	if (isDeleted) {
		json.NewEncoder(w).Encode(ResponseMsg{
			Message: "Successfully deleted Product",
		})
	}
}

func (app *App) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("handleCreateOrder")
	var body OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ResponseMsg := ResponseMsg{
			Message: "Something went wrong.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	order, err := app.createOrder(body)
	if err != nil {
		ResponseMsg := ResponseMsg{
			Message: "Failed to create order",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (app *App) handleGetOrderList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Handle List Orders")

	orders := app.listOrders()
	json.NewEncoder(w).Encode(orders)
}

func (app *App) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("handleGetOrder")

	stringID := r.PathValue("id")

	id, err := strconv.Atoi(stringID)
	if err != nil {
		ResponseMsg := ResponseMsg{
			Message: "Invalid ID.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	order, err := app.getOrder(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			ResponseMsg := ResponseMsg{
				Message: "Order does not exist.",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseMsg)
			return
		}

		ResponseMsg := ResponseMsg{
			Message: "Something went wrong.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func (app *App) handleOrderComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("handleOrderComplete")
	stringID := r.PathValue("id")

	id, err := strconv.Atoi(stringID)
	if err != nil {
		ResponseMsg := ResponseMsg{
			Message: "Invalid ID.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}


	order, err := app.completeOrder(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			ResponseMsg := ResponseMsg{
				Message: "Order does not exist.",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseMsg)
			return
		}

		ResponseMsg := ResponseMsg{
			Message: "Something went wrong.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (app *App) handleOrderCancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("handleOrderCancel")
	stringID := r.PathValue("id")

	id, err := strconv.Atoi(stringID)
	if err != nil {
		ResponseMsg := ResponseMsg{
			Message: "Invalid ID.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}


	order, err := app.cancelOrder(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			ResponseMsg := ResponseMsg{
				Message: "Order does not exist.",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseMsg)
			return
		}

		ResponseMsg := ResponseMsg{
			Message: "Something went wrong.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg)
		return
	}

	json.NewEncoder(w).Encode(order)
}
