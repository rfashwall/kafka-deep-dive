package model

import "github.com/google/uuid"

type Order struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Products []Product `json:"products"`
	Customer Customer  `json:"customer"`
}

type Product struct {
	ProductCode string `json:"productCode"`
	Quantity    int    `json:"quantity"`
}

type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
}

type Customer struct {
	FirstName       string  `json:"firstName"`
	LastName        string  `json:"lastName"`
	EmailAddress    string  `json:"emailAddress"`
	ShippingAddress Address `json:"shippingAddress"`
}
