package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)


type ProductRequest struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type ResponseMsg struct {
	Message string `json:"message"`
}

type App struct {
	pool          *pgxpool.Pool
}

var ErrDBQueryIssue = errors.New("db query issue")
func main() {
	app := App{}

	if err := app.setupDBConnection(false); err != nil {
		log.Fatal(err)
	}
	defer app.pool.Close()

	app.handleRequest()
}

func (app *App) setupDBConnection(isTest bool) error {
	// Load the .env file from the current directory
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("load .env: %w", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if isTest {
		databaseURL = os.Getenv("DATABASE_URL_TEST")
	}

	conn, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return errors.New("Unable to connect to database")
	}
	app.pool = conn
	if err := app.pool.Ping(context.Background()); err != nil {
		app.pool.Close()
		app.pool = nil
		return fmt.Errorf("ping database: %w", err)
	}	

	return nil
}

func (app *App) handleRequest() {
	mux := app.routes()
	fmt.Println("Server running at localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (app *App) routes() http.Handler {
	mux := http.NewServeMux()

	// Product routes
	mux.HandleFunc("GET /products", app.handleListProducts)
	mux.HandleFunc("POST /products", app.handleCreateProduct)

	return mux
}

func (app *App) handleListProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List Products")
	w.Header().Set("Content-Type", "application/json")
	products, err := app.listProducts(r.Context())
	if err != nil {
		if errors.Is(err, ErrDBQueryIssue) {
			ResponseMsg := ResponseMsg{
				Message: "Internal Server Error",
			}
			w.WriteHeader(http.StatusInternalServerError)
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
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

	product, err := app.addProduct(r.Context(), body)
	if err != nil {
		if errors.Is(err, ErrDBQueryIssue) {
			ResponseMsg := ResponseMsg{
				Message: "Internal Server Error",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ResponseMsg)
			return	
		}

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

