package calc_test

import (
	"testing"

	"github.com/mi-wada/go_playground/calc"
)

func TestFastSum(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	got := calc.FastSum(s)
	want := 15
	if got != want {
		t.Errorf("FastSum(%v) = %d, want %d", s, got, want)
	}
}

func BenchmarkFastSum(b *testing.B) {
	n := 10_000
	s := make([]int, n)
	for i := range n {
		s[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		calc.FastSum(s)
	}
}

func TestSlowSum(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	got := calc.SlowSum(s)
	want := 15
	if got != want {
		t.Errorf("FastSum(%v) = %d, want %d", s, got, want)
	}
}

func BenchmarkSlowSum(b *testing.B) {
	n := 10_000
	s := make([]int, n)
	for i := range n {
		s[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		calc.SlowSum(s)
	}
}
