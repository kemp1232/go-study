package main

import (
	"encoding/json"
	"fmt"
	"os"
)
type Product struct {
	ID int	`json:"id"`
	Name string	`json:"name"`
	Category string `json:"category"`
	Price float64	`json:"price"`
	Quantity int	`json:"quantity"`
}

type App struct {
	products map[int]Product
	nextProductID int
}

func main() {

	app := App {
		products: map[int]Product{},
		nextProductID: 1,
	}

	if err := app.loadProducts(); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No Inventory file found. Starting Empty")
		} else {
			fmt.Println("Failed to load products: ", err)
		}
	}

	app.openMenu()
}

func (app *App) loadProducts() error {
	jsonFile, err := os.ReadFile("products.json")
	if err != nil {
		return err
	}

	var products []Product
	err = json.Unmarshal(jsonFile, &products)

	if err != nil {
		return err
	}

	fmt.Println("Successfully opened file!")

	highestID := 0
	for _, product := range products {
		app.products[product.ID] = product

		if product.ID > highestID {
			highestID = product.ID
		}
	}

	app.nextProductID = highestID + 1
	return nil
}

func (app *App) openMenu() {
	for {
		fmt.Println()
		fmt.Println()
		fmt.Println("==== Inventory Manager V2 ====")
		fmt.Println()
		fmt.Println("1. Add Product")
		fmt.Println("2. Update Stock")
		fmt.Println("3. Remove Product")
		fmt.Println("4. List Products")
		fmt.Println("5. Exit")

		fmt.Print("Select an Option: ")
		var choice int
		fmt.Scan(&choice)

		switch choice {
			case 1:
				app.addProduct()
			case 2:
				app.updateStock()
			case 3:
				app.removeProduct()
			case 4:
				app.listProducts()
			case 5:
				fmt.Println("Exit")
				return
			default:
				fmt.Println("Invalid Option. Please Select again.")
		}
	}
}
func (app *App) addProduct() {
	fmt.Println("==== Add Product ====")

	var (
		name string
		category string
		price float64
		quantity int
	)

	fmt.Print("Name: ")
	fmt.Scan(&name)
	fmt.Println()

	fmt.Print("Category: ")
	fmt.Scan(&category)

	for {
		fmt.Print("Price: ")
		fmt.Scan(&price)
		fmt.Println()

		if price > 0 {
			break
		}

		fmt.Println("Price must be greater than 0.")
	}

	for {
		fmt.Print("Quantity: ")
		fmt.Scan(&quantity)
		fmt.Println()

		if quantity >= 0 {
			break
		}

		fmt.Println("Quantity cannot be niggative.")
	}

	id := app.nextProductID

	newProduct := Product{
		ID: id,
		Name: name,
		Category: category,
		Price: price,
		Quantity: quantity,
	}

	app.products[id] = newProduct
	app.nextProductID++

	err := app.saveProducts()
	if err != nil {
		fmt.Println("Error saving Products: ", err)
	}
}

func (app *App) updateStock() {
	fmt.Println("==== Update Product Stock ====")

	var (
		productID int
		quantityChange int
		selectedProduct Product
	)

	for {
		fmt.Print("Product ID: ")
		fmt.Scan(&productID)
		product, exists := app.products[productID]
		if exists {
			selectedProduct = product
			break
		}

		fmt.Println("Product does not exist.")
	}

	for {
		fmt.Print("Quantity Change: ")
		fmt.Scan(&quantityChange)

		newQuantity := selectedProduct.Quantity + quantityChange
		if newQuantity >= 0 {
			selectedProduct.Quantity = newQuantity
			break
		}

		fmt.Println("Quantity cannot be niggative")
	}

	app.products[productID] = selectedProduct
	err := app.saveProducts()
	if err != nil {
		fmt.Println("Error saving Products: ", err)
	}
}

func (app *App) removeProduct() {
	var (
		productID int
	)

	for {
		fmt.Print("Product ID: ")
		fmt.Scan(&productID)
		_, exists := app.products[productID]
		if exists{
			break
		}

		fmt.Println("Product does not exist.")
	}

	delete(app.products, productID)
	err := app.saveProducts()
	if err != nil {
		fmt.Println("Error saving Products: ", err)
	}
}

func (app *App) listProducts() {
	fmt.Println("==== List Products ====")

	for _, product := range app.products {
		fmt.Printf("%v | %v | %v | $%v | %v Pcs \n", product.ID, product.Name, product.Category, product.Price, product.Quantity)
	}
}

func (app *App) saveProducts() error {
	var products []Product

	for _, product := range app.products {
		products = append(products, product)
	}



	jsonString, err := json.MarshalIndent(products, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile("products.json", jsonString, 0644)
	if err != nil {
		return err
	}

	return nil
}

//

// open menu
// add product
// update stock
// remove stprodock
// save to products. json after every update
// 