package main

import (
	"errors"
	"fmt"
)

type Order struct {
	ID     int
	Items  []OrderItem
	Total  float64
	Status string
}

type OrderItem struct {
	ProductID int
	Quantity  int
}

func (app *App) getOrderByID(orderID int) (Order, error) {
	order, exists := app.orders[orderID]
	if !exists {
		return Order{}, errors.New("Order does not exist.")
	}

	return order, nil
}

func (app *App) createOrder(
	orderItems []OrderItem,
) (Order, error) {
	/*
		What is the total?
		What ID should this order receive?
		What status should it have?
	*/


	if len(orderItems) < 1 {
		return Order{}, errors.New("Order does not have items")
	}

	var total float64
	productIds := map[int]int{}
	for _, orderItem := range orderItems {
		product, err := app.getProductByID(orderItem.ProductID)

		// Does each product exist?
		if err != nil {
			errMsg := fmt.Sprintf("Product ID: %v cannot be found", orderItem.ProductID)
			return Order{}, errors.New(errMsg)
		}

		// Is quantity > 0?
		if orderItem.Quantity <= 0 {
			return Order{}, errors.New("Quantity must be greater than 0")
		}

		// Is enough stock available?
		if product.Stock-orderItem.Quantity < 0 {
			return Order{}, errors.New("Not enough stock for the quantity inputted.")
		}

		_, exists := productIds[orderItem.ProductID]
		if exists {
			return Order{}, errors.New("Duplicate Product ID")
		}

		productIds[orderItem.ProductID] = orderItem.ProductID

		total += product.Price * float64(orderItem.Quantity)
	}

	id := app.nextOrderID
	newOrder := Order{
		ID:     id,
		Items:  orderItems,
		Total:  total,
		Status: "pending",
	}
	app.orders[id] = newOrder
	app.nextOrderID++

	return newOrder, nil
}

func (app *App) calculateOrderTotal(orderID int) (float64, error) {
	var (
		total float64
	)

	order, err := app.getOrderByID(orderID)
	if err != nil {
		return 0, err
	}

	for _, orderItem := range order.Items {
		product, exists := app.products[orderItem.ProductID]
		if !exists {
			errMsg := fmt.Sprintf("Product ID: %v does not exist:", orderItem.ProductID)
			return 0, errors.New(errMsg)
		}

		total += product.Price * float64(orderItem.Quantity)
	}

	return total, nil
}

func (app *App) listOrders() {
	for _, order := range app.orders {
		fmt.Printf("%v %v items $%v %v \n", order.ID, len(order.Items), order.Total, order.Status)

		for _, orderItem := range order.Items {
			fmt.Println("= Order Items =")
			fmt.Printf("Product ID: %v %v Pcs | Product Name: %v Product Price Per Piece: $%v \n", orderItem.ProductID, orderItem.Quantity, app.products[orderItem.ProductID].Name, app.products[orderItem.ProductID].Price)
		}
	}
}

func (app *App) completeOrder(orderID int) (Order, error) {
	order, err := app.getOrderByID(orderID)
	if err != nil {
		return Order{}, err
	}

	if order.Status != "pending" {
		errMsg := fmt.Sprintf("Cannot complete order, order status is %v", order.Status)
		return Order{}, errors.New(errMsg)
	}

	if len(order.Items) < 1 {
		return Order{}, errors.New("Order does not have items")
	}

	productIds := map[int]int{}
	for _, orderItem := range order.Items {
		product, exists := app.products[orderItem.ProductID]
		if !exists {
			errMsg := fmt.Sprintf("Product ID: %v cannot be found", orderItem.ProductID)
			return Order{}, errors.New(errMsg)
		}

		remainingStock := product.Stock - orderItem.Quantity
		if remainingStock < 0 {
			errMsg := fmt.Sprintf("Not enough stocks. You ordered: %v pcs, Remaining stock is: %v", orderItem.Quantity, product.Stock)
			return Order{}, errors.New(errMsg)
		}

		_, idExists := productIds[orderItem.ProductID]
		if idExists {
			return Order{}, errors.New("Invalid Order, Duplicate Product ID")
		}

		productIds[orderItem.ProductID] = orderItem.ProductID
	}

	for _, orderItem := range order.Items {
		product, _ := app.products[orderItem.ProductID]
		remainingStock := product.Stock - orderItem.Quantity
		product.Stock = remainingStock
		app.products[orderItem.ProductID] = product
	}

	order.Status = "completed"
	app.orders[orderID] = order
	return order, nil
}

func (app *App) cancelOrder(orderID int) (Order, error) {
	order, err := app.getOrderByID(orderID)
	if err != nil {
		return Order{}, err
	}

	if order.Status != "pending" {
		errMsg := fmt.Sprintf("Cannot cancel order, order status is %v", order.Status)
		return Order{}, errors.New(errMsg)
	}

	order.Status = "cancelled"
	app.orders[orderID] = order
	return order, nil
}
