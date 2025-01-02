package calc_test

import (
	"testing"

	"github.com/mi-wada/go_playground/calc"
)

func TestAdd(t *testing.T) {
	for _, test := range []struct {
		name           string
		lhs, rhs, want int
	}{
		{"pos + pos", 1, 2, 3},
		{"pos + neg", 1, -2, -1},
		{"neg + pos", -1, 2, 1},
		{"neg + neg", -1, -2, -3},
		{"zero + pos", 0, 2, 2},
		{"pos + zero", 1, 0, 1},
		{"zero + neg", 0, -2, -2},
		{"neg + zero", -1, 0, -1},
		{"zero + zero", 0, 0, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := calc.Add(test.lhs, test.rhs)
			if got != test.want {
				t.Errorf("calc.Add(%d, %d) = %d, want %d", test.lhs, test.rhs, got, test.want)
			}
		})
	}
}
