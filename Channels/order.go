package main

import "fmt"

type order struct {
	ProductCode int
	Quantity float64
	Status orderStatus
}

type invalidOrder struct {
	order order
	err error
}

/* stringer interface, what this will do is when you print the object, using a function from the FMT package, 
it will call the string method and use the result of the string method.*/
func(o order) String() string {
	return fmt.Sprintf("Product code: %v, Quantity: %v, Status: %v\n", o.ProductCode, o.Quantity, orderStatusToText(o.Status))
}

func orderStatusToText(o orderStatus) string {
	switch o {
	case none:
		return "none"
	case new:
		return "new"
	case received:
		return "received"
	case reserved:
		return "reserved"
	case filled:
		return "filled"
	default:
		return "unknown status"
	}
}

type orderStatus int

const (
	// iota is a great way in Go to allow us to have unique values inside of a group of constants.
	none orderStatus= iota
	new
	received
	reserved
	filled
)
