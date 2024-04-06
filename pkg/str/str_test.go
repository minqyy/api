package str

import "testing"

func TestCompleteStringToLength(t *testing.T) {
	// Test case: Empty string
	emptyString := ""
	if result := CompleteStringToLength(emptyString, 5, '*'); result != "*****" {
		t.Errorf("Expected '*****' for empty string, got '%s'", result)
	}

	// Test case: String length is equal to target length
	inputString := "hello"
	if result := CompleteStringToLength(inputString, 5, '*'); result != "hello" {
		t.Errorf("Expected 'hello' for string length equal to target length, got '%s'", result)
	}

	// Test case: String length is greater than target length
	if result := CompleteStringToLength(inputString, 3, '*'); result != "hel" {
		t.Errorf("Expected 'hel' for string length greater than target length, got '%s'", result)
	}

	// Test case: String length is less than target length
	if result := CompleteStringToLength(inputString, 10, '*'); result != "hello*****" {
		t.Errorf("Expected 'hello*****' for string length less than target length, got '%s'", result)
	}

	// Test case: Target length is negative
	if result := CompleteStringToLength(inputString, -5, '*'); result != "" {
		t.Errorf("Expected empty string for negative target length, got '%s'", result)
	}

	// Test case: Char is a space
	if result := CompleteStringToLength(inputString, 10, ' '); result != "hello     " {
		t.Errorf("Expected 'hello     ' for char being space, got '%s'", result)
	}

	// Test case: Char is a special character
	if result := CompleteStringToLength(inputString, 10, '@'); result != "hello@@@@@" {
		t.Errorf("Expected 'hello@@@@@' for char being '@', got '%s'", result)
	}
}
