package slice

import "testing"

func TestContains(t *testing.T) {
	// Test case: Empty slice
	var emptySlice []int
	if Contains(emptySlice, 1) {
		t.Error("Expected false for empty slice, got true")
	}

	// Test case: Slice with one element that matches
	sliceWithOne := []string{"apple"}
	if !Contains(sliceWithOne, "apple") {
		t.Error("Expected true for slice with one element that matches, got false")
	}

	// Test case: Slice with one element that does not match
	if Contains(sliceWithOne, "banana") {
		t.Error("Expected false for slice with one element that does not match, got true")
	}

	// Test case: Slice with multiple elements including the searched one
	sliceWithMultiple := []int{1, 2, 3, 4, 5}
	if !Contains(sliceWithMultiple, 3) {
		t.Error("Expected true for slice with multiple elements including the searched one, got false")
	}

	// Test case: Slice with multiple elements but without the searched one
	if Contains(sliceWithMultiple, 6) {
		t.Error("Expected false for slice with multiple elements but without the searched one, got true")
	}

	// Test case: Slice with nil value
	var nilSlice []float64
	if Contains(nilSlice, 0.0) {
		t.Error("Expected false for nil slice, got true")
	}

	// Test case: Slice with non-nil value but empty interface type
	interfaceSlice := make([]interface{}, 0)
	if Contains(interfaceSlice, "test") {
		t.Error("Expected false for slice with empty interface type, got true")
	}

	// Test case: Slice with nil element
	sliceWithNil := []interface{}{nil, 1, 2, 3}
	if !Contains(sliceWithNil, nil) {
		t.Error("Expected true for slice with nil element, got false")
	}
}
