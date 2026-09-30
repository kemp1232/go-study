package main

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
