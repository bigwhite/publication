package main

import "testing"

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2, 3) = %d; want 5", result)
	}

	result = Add(0, 5)
	if result != 5 {
		t.Errorf("Add(0, 5) = %d; want 5", result)
	}
}

func TestIsPositive(t *testing.T) {
	result := IsPositive(10)
	if !result {
		t.Error("IsPositive(10) = false; want true")
	}

	result = IsPositive(-5)
	if result {
		t.Error("IsPositive(-5) = true; want false")
	}

	result = IsPositive(0)
	if result {
		t.Error("IsPositive(0) = true; want false")
	}
}
