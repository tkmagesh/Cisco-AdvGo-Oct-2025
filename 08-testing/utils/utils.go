package utils

import "math"

func IsPrime(no int) bool {
	for i := float64(2); i <= math.Sqrt(float64(no)); i++ {
		if no%int(i) == 0 {
			return false
		}
	}
	return true
}

func GeneratePrimes(start, end int) []int {
	primes := make([]int, 0, 100)
	// var primes []int
	for no := start; no <= end; no++ {
		if IsPrime(no) {
			primes = append(primes, no)
		}
	}
	return primes
}
