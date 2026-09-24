package main

import (
	"errors"
	"fmt"
	"strings"
)

type Product struct {
	ID int
	Name string
	Category string
	Price float64
	Quantity int
}

type App struct {
	products map[int]Product
	nextProductID int
}



func main() {
	app := App{
		products: map[int]Product{},
		nextProductID: 1,
	}

	app.openMenu()
}

func (app *App) openMenu() {
	for {
		fmt.Println("")
		fmt.Println("")
		fmt.Println("=== Inventory Manager ===")
		fmt.Println("")
		fmt.Println("1. Add Product")
		fmt.Println("2. List Products")
		fmt.Println("3. Find Product")
		fmt.Println("4. Update Stock")
		fmt.Println("5. Remove Product")
		fmt.Println("6. Calculate Inventory Value")
		fmt.Println("7. List Products By Category")
		fmt.Println("8. Exit")
	
		fmt.Print("Select an Option: ")
		var choice int
		fmt.Scan(&choice)

		switch choice {
			case 1:
				app.addProduct()
			case 2:
				app.listProducts()
			case 3:
				app.findProduct()
			case 4:
				app.updateStock()
			case 5:
				app.removeProduct()
			case 6:
				app.calculateInventoryValue()
			case 7:
				app.listProductsByCategory()
			case 8:
				fmt.Println("Exiting app...")
				return
		}
	}
}

func (app *App) findProductByID(productID int) (Product, error) {
	product, exists := app.products[productID]

	if !exists {
		return Product{}, errors.New("product cannot be found")
	}

	return product, nil
}

func (app *App) addProduct() {
	fmt.Println("=== Add Product ===")

	var (
		name string
		category string
		price float64
		quantity int
	)

	fmt.Print("Name: ")
	fmt.Scan(&name)
	fmt.Println("")

	fmt.Print("Category: ")
	fmt.Scan(&category)
	fmt.Println("")

	for {
		fmt.Print("Price: ")
		fmt.Scan(&price)
		fmt.Println("")
		if price > 0 {
			break
		}

		fmt.Println("Price must be greater than 0.")
	}

	for {
		fmt.Print("Quantity: ")
		fmt.Scan(&quantity)
		fmt.Println("")
		if (quantity >= 0) {
			break
		}

		fmt.Println("Quantity cannot be negative.")
	}

	ID := app.nextProductID
	newProduct := Product{
		ID: ID,
		Name: name,
		Category: category,
		Price: price,
		Quantity: quantity,
	}

	app.products[ID] = newProduct
	app.nextProductID++
}


func (app *App) listProducts() {
	fmt.Println("=== List Products ===")
	for _, product := range app.products {
		fmt.Printf(
			"%v %v %v $%v %v Pcs \n",
			product.ID,
			product.Name,
			product.Category,
			product.Price,
			product.Quantity,
		)
	}
}



func (app *App) findProduct() {
	fmt.Println("=== Find Product ===")

	var productID int
	fmt.Print("Product ID: ")
	fmt.Scan(&productID)
	fmt.Println("")


	product, err := app.findProductByID(productID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ID: ", product.ID)
	fmt.Println("Name: ", product.Name)
	fmt.Println("Category: ", product.Category)
	fmt.Println("Price: ", product.Price)
	fmt.Println("Quantity: ", product.Quantity)
}

func (app *App) updateStock() {
	var productID int
	var quantityChange int
	var selectedProduct Product
	var updatedQuantity int
	for {
		fmt.Print("Product ID: ")
		fmt.Scan(&productID)
		fmt.Println("")

		fetchedProduct, err := app.findProductByID(productID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		selectedProduct = fetchedProduct
		break
	}

	for {
		fmt.Print("Quantity Change: ")
		fmt.Scan(&quantityChange)
		fmt.Println("")

		updatedQuantity = selectedProduct.Quantity + quantityChange
		isInvalidQuantityChange := updatedQuantity < 0
		if !isInvalidQuantityChange {
			break
		} 

		fmt.Println("Insufficient stock.")
	}


	selectedProduct.Quantity = updatedQuantity
	app.products[productID] = selectedProduct
}

func (app *App) removeProduct() {
	var productID int
	for {
		fmt.Print("Product ID: ")
		fmt.Scan(&productID)
		fmt.Println("")
		_, err := app.findProductByID(productID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		
		break
	}

	delete(app.products, productID)
}

func (app *App) calculateInventoryValue() {
	var totalInventoryValue float64

	for _, product := range app.products {
		
		totalInventoryValue += (product.Price * float64(product.Quantity))
	}

	fmt.Printf("Total Inventory Value is: $%v \n", totalInventoryValue)
}

func (app *App) listProductsByCategory() {
	var productCategory string
	fmt.Print("Product Category: ")
	fmt.Scan(&productCategory)

	for _, product := range app.products {
		if strings.EqualFold(product.Category, productCategory) {
			fmt.Printf(
				"%v %v %v $%v %v Pcs \n",
				product.ID,
				product.Name,
				product.Category,
				product.Price,
				product.Quantity,
			)
		}
	}
}
