package main

import (
	"demo/Testing_in_Go/Writing_a_unit_test/user"
	"log"
	"net/http"
)

func main() {
	const address = ":3000"

	http.HandleFunc("/users/", user.Handler)

	log.Fatal(http.ListenAndServe(address, nil))
}
