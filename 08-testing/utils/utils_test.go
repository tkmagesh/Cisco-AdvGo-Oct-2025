package utils

import (
	"fmt"
	"testing"
)

func TestIsPrime_19(t *testing.T) {
	// Arrange
	no := 19
	expectedResult := true

	// to force the test to fail
	// expectedResult := false

	// Act
	actualResult := IsPrime(no)

	t.Log("Verifying if 19 is prime\n")

	// Assert
	if actualResult != expectedResult {
		/*
			t.Logf("expected = %v, but actual = %v\n", expectedResult, actualResult)
			t.Fail()
		*/
		t.Errorf("expected = %v, but actual = %v\n", expectedResult, actualResult)
	}
}

func TestIsPrime_27(t *testing.T) {

	// skipping this test
	// t.Skip("skipping IsPrime(27)")

	// Arrange
	no := 27
	expectedResult := false

	// Act
	actualResult := IsPrime(no)

	t.Log("Verifying if 27 is prime\n")

	// Assert
	if actualResult != expectedResult {
		/*
			t.Logf("expected = %v, but actual = %v\n", expectedResult, actualResult)
			t.Fail()
		*/
		t.Errorf("expected = %v, but actual = %v\n", expectedResult, actualResult)
	}
}

func TestDummy(t *testing.T) {

}

func TestPrimes(t *testing.T) {
	testData := []struct {
		no       int
		expected bool
	}{
		{no: 11, expected: true},
		{no: 13, expected: true},
		{no: 19, expected: true},
		{no: 17, expected: true},
		{no: 29, expected: true},
	}
	for _, td := range testData {
		t.Run(fmt.Sprintf("IsPrime(%d)", td.no), func(t *testing.T) {
			actual := IsPrime(td.no)

			if actual != td.expected {
				t.Errorf("expected = %v, but actual = %v\n", td.expected, actual)
			}
		})
	}
}

// Benchmarking
func BenchmarkGeneratePrimes(b *testing.B) {
	for b.Loop() {
		GeneratePrimes(2, 100)
	}
}
