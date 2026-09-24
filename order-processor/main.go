package main

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


type Product struct {
	ID int	
	Name string	
	Price float64	
	Stock int	
}

type Order struct {
	ID int	
	Items []OrderItem
	Total float64	
	Status int	
}

type OrderItem struct {
	ProductID int	
	Quantity int	
}


type App struct {
	products map[int]Product
	orders map[int]Order
	nextProductID int
	nextOrderID int
}

func main() {

}