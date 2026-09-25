package main

import (
	"fmt"
)

/*
Product
- ID
- Name
- Price
- Stock

OrderItem
- ProductID
- Quantity

Order
- ID
- Items
- Total
- Status


pending
completed
cancelled
*/

type App struct {
	products map[int]Product
	orders map[int]Order
	nextProductID int
	nextOrderID int
}

func main() {
	app := App {
		products: map[int]Product{},
		orders: map[int]Order{},
		nextProductID: 5,
		nextOrderID: 1,
	}

	// Default values for trial and error
	app.products[1] = Product{
		ID: 1,
		Name: "CBR 650R",
		Price: 8000,
		Stock: 20,
	}

	app.products[2] = Product{
		ID: 2,
		Name: "Speed 400",
		Price: 3000,
		Stock: 10,
	}

	app.products[3] = Product{
		ID: 3,
		Name: "Dominar 400 UG",
		Price: 2000,
		Stock: 30,
	}

	app.products[4] = Product{
		ID: 4,
		Name: "Z650RS",
		Price: 3500,
		Stock: 10,
	}

	app.openMenu()
}

func (app *App) openMenu() {
	for {
		fmt.Println()
		fmt.Println()
		fmt.Println("==== Inventory Processor ====")
		fmt.Println()
		fmt.Println("1. Add Product")
		fmt.Println("2. List Products")
		fmt.Println("3. Create Order")
		fmt.Println("4. Calculate Order Total")
		fmt.Println("5. Complete Order")
		fmt.Println("6. Cancel Order")
		fmt.Println("7. List Orders")
		fmt.Println("8. Exit")

		fmt.Print("Select an Option: ")
		var choice int
		fmt.Scan(&choice)

		switch choice {
			case 1:
				app.handleAddProduct()
			case 2:
				app.handleListProducts()
			case 3:
				app.handleCreateOrder()
			case 4:
				app.handleCalculateOrderTotal()
			case 5:
				app.handleCompleteOrder()
			case 6:
				app.handleCancelOrder()
			case 7:
				app.handleListOrders()
			case 8:
				fmt.Println("Exit")
				return
			default:
				fmt.Println("Invalid Option. Please Select again.")
		}
	}
}


func (app *App) handleAddProduct() {
	fmt.Println("==== Add Product ====")

	var (
		name string
		price float64
		stock int
	)

	fmt.Print("Name: ")
	fmt.Scan(&name)
	fmt.Println()

	for {
		fmt.Print("Price: ")
		fmt.Scan(&price)

		if price > 0 {
			break
		}

		fmt.Println("Price must be greater than 0.")
	}

	for {
		fmt.Print("Stock: ")
		fmt.Scan(&stock)
		fmt.Println()

		if stock >= 0 {
			break
		}

		fmt.Println("Stock must not be negative.")
	}

	app.addProduct(name, price, stock)
}

func (app *App) handleListProducts() {
	fmt.Println("==== List Product ====")
	app.listProducts()
}

func (app *App) handleCreateOrder() {
	fmt.Println("==== Create Order ====")
	fmt.Println("==== Available Products ====")
	app.listProducts()

	var (
		orders []OrderItem
		// total float64
	)

	OuterLoop: 
		for {
			var (
				quantity int
				productID int
			)

			fmt.Print("Product ID: ")
			fmt.Scan(&productID)
			fmt.Println()

			fmt.Print("Quantity: ")
			fmt.Scan(&quantity)
			fmt.Println()


			orderItem := OrderItem{
				ProductID: productID,
				Quantity: quantity,
			}
			orders = append(orders, orderItem)

			var isAddAnotherItem string
			for {
				fmt.Print("Would you like to add another item? (y/n): ")
				fmt.Scan(&isAddAnotherItem)
				fmt.Println()
		
				if isAddAnotherItem == "n" {
					break OuterLoop
				}
		
				if isAddAnotherItem == "y" {
					break
				}

				fmt.Println("Invalid input.")
			}
		}

	order, err := app.createOrder(orders)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Created order:", order.ID)
}


func (app *App) handleCalculateOrderTotal() {
	fmt.Println("==== Calculate Order Total ====")
	var (
		orderID int
	)

	fmt.Print("Order ID: ")
	fmt.Scan(&orderID)
	fmt.Println()

	total, err := app.calculateOrderTotal(orderID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Order with ID of: %v, has a total of $%v \n", orderID, total)
}

func (app *App) handleCompleteOrder() {
	fmt.Println("==== Complete Order ====")
	var (
		orderID int
	)
	fmt.Print("Order ID: ")
	fmt.Scan(&orderID)
	fmt.Println()

	app.setOrderStatus(orderID, "completed")
}

func (app *App) handleCancelOrder() {
	fmt.Println("==== Cancel Order ====")
	var (
		orderID int
	)
	fmt.Print("Order ID: ")
	fmt.Scan(&orderID)
	fmt.Println()

	app.setOrderStatus(orderID, "cancelled")
}

func (app *App) handleListOrders() {
	fmt.Println("==== List Orders ====")
	app.listOrders()
}