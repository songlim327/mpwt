package core

import "errors"

// CalculatePaneSize takes an integer (number of panes) and returns a slice of float64
func CalculatePaneSize(n int) ([]float64, error) {
	if n < 1 {
		return nil, errors.New("length must be greater than 0")
	}

	results := make([]float64, n)

	for i := 0; i < n; i++ {
		num := float64(n - i)
		denom := float64(n - i + 1)
		results[i] = num / denom
	}

	return results, nil
}
