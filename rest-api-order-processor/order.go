package main

import (
	"errors"
	"fmt"
)

type Order struct {
	ID     int         `json:"id"`
	Items  []OrderItem `json:"items"`
	Total  float64     `json:"total"`
	Status string      `json:"status"`
}

type OrderItem struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}

func (app *App) createOrder(order OrderRequest) (Order, error) {

	if len(order.Items) < 1 {
		return Order{}, errors.New("Order does not have items")
	}

	var total float64
	productIds := map[int]int{}
	for _, orderItem := range order.Items {
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
		if product.Quantity-orderItem.Quantity < 0 {
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
		Items:  order.Items,
		Total:  total,
		Status: "pending",
	}
	app.orders[id] = newOrder
	app.nextOrderID++

	return newOrder, nil
}


func (app *App) listOrders() map[int]Order {
	return app.orders
}

func (app *App) getOrder(orderID int) (Order, error) {
	order, exists := app.orders[orderID]
	if !exists {
		return Order{}, ErrNotFound
	}

	return order, nil
}

func (app *App) completeOrder(orderID int) (Order, error) {
	order, err := app.getOrder(orderID)
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

		remainingStock := product.Quantity - orderItem.Quantity
		if remainingStock < 0 {
			errMsg := fmt.Sprintf("Not enough stocks. You ordered: %v pcs, Remaining stock is: %v", orderItem.Quantity, product.Quantity)
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
		remainingStock := product.Quantity - orderItem.Quantity
		product.Quantity = remainingStock
		app.products[orderItem.ProductID] = product
	}

	order.Status = "completed"
	app.orders[orderID] = order
	return order, nil
}

func (app *App) cancelOrder(orderID int) (Order, error) {
	order, err := app.getOrder(orderID)
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