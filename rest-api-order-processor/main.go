package main

import (
	"fmt"
	"net/http"
)

type App struct {
	products      map[int]Product
	orders        map[int]Order
	nextProductID int
	nextOrderID   int
}

func (app *App) main() {
	app.handleRequest()
}

func (app *App) handleRequest() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", handleGetProduct)
	mux.HandleFunc("POST /products", handleCreateProduct)
	mux.HandleFunc("GET /products/{id}", handleGetProduct)
	mux.HandleFunc("PUT /products/{id}", handleUpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", handleDeleteProduct)

	/*

		POST   /orders                 Create order
		GET    /orders                 List orders
		GET    /orders/{id}            Get order
		POST   /orders/{id}/complete   Complete order
		POST   /orders/{id}/cancel     Cancel order
	*/
}

func handleListProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleListProducts")
}
func handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleCreateProduct")
}
func handleGetProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleGetProduct")
}
func handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleUpdateProduct")
}
func handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleDeleteProduct")
}
