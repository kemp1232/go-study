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

func TestgetOrderByIDNotExist(t *testing.T) {
	app := App{}

	orderID, err := app.getOrderByID(1)
	if err == nil {
		t.Error("Expected an error for non-existing order id, got nil")
	}

	if orderID != 0 {
		t.Errorf("Expected orderID is 0 got %v", orderID)
	}
}

func TestgetOrderByID(t *testing.T) {
	app := App{
		orders: map[int]Order{},
	}

	app.orders[123] = Order{
		ID: 123,
		Items: []OrderItem{},
		Total: 8000,
		Status: "pending",
	}

	orderID, err := app.getOrderByID(123)

	if err != nil {
		t.Errorf("Did not Expect an error, got %v", err)
	}

	if orderID != 123 {
		t.Errorf("Expected orderID is 123 got %v", orderID)
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