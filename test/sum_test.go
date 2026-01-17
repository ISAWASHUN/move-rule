package test

import "testing"

func TestSum(t *testing.T) {
	result := Sum(1, 2)
	if result != 3 {
		t.Errorf("Sum(1, 2) = %d; want 3", result)
	}
}
