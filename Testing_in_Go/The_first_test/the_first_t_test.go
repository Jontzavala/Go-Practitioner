package main

import "testing"

func TestAdd(t *testing.T) {
	l, r := 1, 2
	expect := 3

	got := add(l, r)

	if expect != got {
		// allows us to report error in the test and f means formatted print statement.
		t.Errorf("Expected %v when adding %v and %v. Got %v\n", expect, l, r, got)
	}
}
