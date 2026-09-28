package main

import (
	"testing"
)

/*

func TestHelloName(t *testing.T) {
	name := "Gladys"
	want := regexp.MustCompile(`\b`+name+`\b`)
	msg, err := Hello("Gladys")s
	if !want.MatchString(msg) || err != nil {
		t.Errorf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Errorf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}
*/

func TestGetOrderByIDNotExist(t *testing.T) {
	app := App{}

	_, err := app.getOrderByID(1)
	if err == nil {
		t.Error("Expected an error for non-existing order id, got nil")
	}
}

func TestGetOrderByID(t *testing.T) {
	app := App{
		orders: map[int]Order{},
	}

	app.orders[123] = Order{
		ID: 123,
		Items: []OrderItem{},
		Total: 8000,
		Status: "pending",
	}

	order, err := app.getOrderByID(123)

	if err != nil {
		t.Errorf("Did not Expect an error, got %v", err)
	}

	if order.ID != 123 {
		t.Errorf("Expected orderID is 123 got %v", order.ID)
	}
}


func TestCalculateOrderTotalMissingProduct(t *testing.T) {
	app := App{
		orders: map[int]Order{},
	}

	orderItems := append(
		app.orders[123].Items,
		OrderItem{
			ProductID: 1,
			Quantity: 10,
		},
	)

	app.orders[123] = Order{
		ID: 123,
		Items: orderItems,
		Total: 8000,
		Status: "pending",
	}

	total, err := app.calculateOrderTotal(123)

	if err == nil {
		t.Errorf("Expected an error for non-existing product id, got nil")
	}

	if total != 0 {
		t.Errorf("Expected total is 0 got %v", total)
	}
}

func TestCalculateOrderTotal(t *testing.T) {
	app := App{
		orders: map[int]Order{},
		products: map[int]Product{},
	}

	orderItems := append(
		app.orders[123].Items,
		OrderItem{
			ProductID: 1,
			Quantity: 10,
		},
	)

	app.orders[123] = Order{
		ID: 123,
		Items: orderItems,
		Total: 0,
		Status: "pending",
	}

	app.products[1] = Product{
		ID: 1,
		Name: "CBR 650R",
		Price: 8000,
		Stock: 20,
	}

	total, err := app.calculateOrderTotal(123)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if total != 80000 {
		t.Errorf("Expected total is 80000 got %v", total)
	}
}


func TestCompletedOrder(t *testing.T) {
	tests := []struct {
		name string
		orderID int
		orderItemID int
		productID int
		quantity int
		stock int
		orderStatus string
		wantErr bool
	}{
		{"Invalid ID",         111, 1, 23, 10, 10, "pending", true},
		{"Invalid Status",     123, 1, 23, 10, 10, "cancelled", true},
		{"Missing Product ID", 123, 1, 55, 10, 10, "pending", true},
		{"Not enough Stocks",  123, 1, 23, 10, 9, "pending", true},
		{"Complete Order",     123, 1, 23, 10, 10, "pending", false},
	}

	app := App{
		orders: map[int]Order{},
		products: map[int]Product{},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			app.products[23] = Product{
				ID: 23,
				Name: "CBR 650R",
				Price: 8000,
				Stock: test.stock,
			}

			orderItems := append(
				app.orders[123].Items,
				OrderItem{
					ProductID: test.productID,
					Quantity: test.quantity,
				},
			)

			app.orders[123] = Order{
				ID: 123,
				Items: orderItems,
				Total: 8000,
				Status: test.orderStatus,
			}


			_, err := app.completeOrder(test.orderID)

			if test.wantErr && err == nil {
				t.Error("Expected an error but got nil")
			}
		})
	}
}

func TestCancelOrder(t *testing.T) {
	tests := []struct {
		name string
		orderID int
		orderStatus string
		wantErr bool
	}{
		{"Invalid ID",     111, "pending", true},
		{"Invalid Status", 123, "cancelled", true},
		{"Cancel Order",   123, "pending", false},
	}

	app := App{
		orders: map[int]Order{},
		products: map[int]Product{},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			app.products[23] = Product{
				ID: 23,
				Name: "CBR 650R",
				Price: 8000,
				Stock: 10,
			}

			orderItems := append(
				app.orders[123].Items,
				OrderItem{
					ProductID: 23,
					Quantity: 10,
				},
			)

			app.orders[123] = Order{
				ID: 123,
				Items: orderItems,
				Total: 8000,
				Status: test.orderStatus,
			}


			_, err := app.cancelOrder(test.orderID)

			if test.wantErr && err == nil {
				t.Error("Expected an error but got nil")
			}
		})
	}
}