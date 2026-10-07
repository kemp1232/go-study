package main

import (
	"fmt"
	"testing"
)

func TestGetOrderNotExist(t *testing.T) {
	app := App{}

	_, err := app.getOrder(1)
	if err == nil {
		t.Error("Expected an error for non-existing order id, got nil")
	}
}

func TestGetOrder(t *testing.T) {
	app := App{
		orders: map[int]Order{},
	}

	app.orders[123] = Order{
		ID:     123,
		Items:  []OrderItem{},
		Total:  8000,
		Status: "pending",
	}

	order, err := app.getOrder(123)

	if err != nil {
		t.Errorf("Did not Expect an error, got %v", err)
	}

	if order.ID != 123 {
		t.Errorf("Expected orderID is 123 got %v", order.ID)
	}
}

func TestCompletedOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderID     int
		orderItemID int
		productID   int
		quantity    int
		prodQuantity       int
		orderStatus string
		wantErr     bool
	}{
		{"Invalid ID", 111, 1, 23, 10, 10, "pending", true},
		{"Invalid Status", 123, 1, 23, 10, 10, "cancelled", true},
		{"Reject Empty Orders", 1, 1, 23, 10, 10, "pending", true},
		{"Duplicate Order", 2, 1, 23, 10, 10, "pending", true},
		{"Missing Product ID", 123, 1, 55, 10, 10, "pending", true},
		{"Not enough Stocks", 123, 1, 23, 10, 9, "pending", true},
		{"Complete Order", 123, 1, 23, 10, 10, "pending", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := App{
				orders:   map[int]Order{},
				products: map[int]Product{},
			}

			app.products[23] = Product{
				ID:    23,
				Name:  "CBR 650R",
				Price: 8000,
				Quantity: test.prodQuantity,
			}

			orderItems := append(
				app.orders[123].Items,
				OrderItem{
					ProductID: test.productID,
					Quantity:  test.quantity,
				},
			)

			app.orders[123] = Order{
				ID:     123,
				Items:  orderItems,
				Total:  8000,
				Status: test.orderStatus,
			}

			app.orders[1] = Order{
				ID:     1,
				Items:  []OrderItem{},
				Total:  8000,
				Status: test.orderStatus,
			}

			//Order with duplicate Orders
			duplicateOrderItems := append(
				app.orders[2].Items,
				OrderItem{
					ProductID: test.productID,
					Quantity:  test.quantity,
				},
			)

			duplicateOrderItems = append(
				duplicateOrderItems,
				OrderItem{
					ProductID: test.productID,
					Quantity:  test.quantity,
				},
			)

			app.orders[2] = Order{
				ID:     2,
				Items:  duplicateOrderItems,
				Total:  8000,
				Status: test.orderStatus,
			}

			order, err := app.completeOrder(test.orderID)

			if test.wantErr && err == nil {
				t.Error("Expected an error but got nil")
			}

			if !test.wantErr && err != nil {
				errMsg := fmt.Sprintf("Did not expected an error but got: %v", err)
				t.Error(errMsg)
			}

			if !test.wantErr && err == nil {
				if order.Status != "completed" {
					t.Error("Order status should have been completed")
				}

				if order.Status == "completed" {
					properProductStock := test.prodQuantity - test.quantity

					if app.products[test.productID].Quantity != properProductStock {
						t.Error("Product Stock was not properly deducted")
					}

					_, errSecondTry := app.completeOrder(test.orderID)

					if errSecondTry == nil {
						t.Error("Expected an error but got nil")
					}
				}
			}
		})
	}
}

func TestCancelOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderID     int
		orderStatus string
		wantErr     bool
	}{
		{"Invalid ID", 111, "pending", true},
		{"Invalid Status", 123, "cancelled", true},
		{"Invalid Status", 123, "completed", true},
		{"Cancel Order", 123, "pending", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := App{
				orders:   map[int]Order{},
				products: map[int]Product{},
			}

			initialStock := 10

			app.products[23] = Product{
				ID:    23,
				Name:  "CBR 650R",
				Price: 8000,
				Quantity: initialStock,
			}

			orderItems := append(
				app.orders[123].Items,
				OrderItem{
					ProductID: 23,
					Quantity:  10,
				},
			)

			app.orders[123] = Order{
				ID:     123,
				Items:  orderItems,
				Total:  8000,
				Status: test.orderStatus,
			}

			order, err := app.cancelOrder(test.orderID)

			if test.wantErr && err == nil {
				t.Error("Expected an error but got nil")
			}

			if !test.wantErr && err != nil {
				errMsg := fmt.Sprintf("Did not expected an error but got: %v", err)
				t.Error(errMsg)
			}

			if !test.wantErr && err == nil {
				if order.Status != "cancelled" {
					t.Error("Order status should have been cancelled")
				}

				if app.products[23].Quantity != initialStock {
					t.Error("Product stock still changed even if the order is cancelled")
				}
			}
		})
	}
}

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name      string
		productID int
		quantity  int
		stock     int
		wantErr   bool
	}{
		{"Invalid Product ID", 123, 10, 10, true},
		{"Insuficient stock", 23, 11, 10, true},
		{"Negative Quantity", 23, -11, 10, true},
		{"Zero Quantity", 23, 0, 10, true},
		{"Create Order", 23, 10, 10, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := App{
				orders:      map[int]Order{},
				products:    map[int]Product{},
				nextOrderID: 1,
			}

			var productPrice float64 = 8000

			app.products[23] = Product{
				ID:    23,
				Name:  "CBR 650R",
				Price: productPrice,
				Quantity: test.stock,
			}

			orderItems := []OrderItem{}
			orderItems = append(orderItems, OrderItem{
				ProductID: test.productID,
				Quantity:  test.quantity,
			})

			orderRequest := OrderRequest{
				Items: orderItems,
			}

			order, err := app.createOrder(orderRequest)

			if test.wantErr && err == nil {
				t.Error("Expected an error but got nil")
			}

			if !test.wantErr && err != nil {
				errMsg := fmt.Sprintf("Did not expected an error but got: %v", err)
				t.Error(errMsg)
			}

			if !test.wantErr && err == nil {
				if order.Status != "pending" {
					t.Error("Order status should have been pending")
				}

				if app.nextOrderID != 2 {
					t.Error("nextOrderID was not properly incremented")
				}

				var total float64
				for _, orderItem := range order.Items {
					total += float64(orderItem.Quantity) * productPrice
				}

				if order.Total != total {
					t.Error("Order total is incorrect")
				}
			}
		})
	}
}

func TestCreateOrderEmpty(t *testing.T) {
	app := App{
		orders:   map[int]Order{},
		products: map[int]Product{},
	}

	orderRequest := OrderRequest{
		Items: []OrderItem{},
	}
	_, err := app.createOrder(orderRequest)

	if err == nil {
		t.Error("expected error for empty order")
	}
}

func TestCreateOrderDuplicateProductID(t *testing.T) {
	app := App{
		orders:   map[int]Order{},
		products: map[int]Product{},
	}

	var productPrice float64 = 8000

	app.products[23] = Product{
		ID:    23,
		Name:  "CBR 650R",
		Price: productPrice,
		Quantity: 10,
	}

	items := []OrderItem{
		{ProductID: 23, Quantity: 2},
		{ProductID: 23, Quantity: 3},
	}

	orderRequest := OrderRequest{
		Items: items,
	}

	_, err := app.createOrder(orderRequest)

	if err == nil {
		t.Error("expected error for empty order")
	}
}
