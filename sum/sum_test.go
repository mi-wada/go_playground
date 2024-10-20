package sum

import (
	"math/rand"
	"testing"
)

const (
	maxSliceSize = 2_000_000
	maxItemValue = 100
)

var (
	tests = []test{
		{"All positive", []int{1, 2, 3, 4, 5}, 15},
		{"All negative", []int{-1, -2, -3, -4, -5}, -15},
		{"Mixed positive and negative", []int{1, -2, 3, -4, 5}, 3},
		{"Empty slice", []int{}, 0},
		{"Single element", []int{10}, 10},
		{"Zeros", []int{0, 0, 0}, 0},
	}
)

type test struct {
	name string
	arg  []int
	want int
}

func TestSum(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sum(tt.arg)
			if result != tt.want {
				t.Errorf("sum(%v) = %d; want %d", tt.arg, result, tt.want)
			}
		})
	}
}

func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sum(generateRandomSlice())
	}
}

func TestSumRec(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sumRec(tt.arg)
			if result != tt.want {
				t.Errorf("sumRec(%v) = %d; want %d", tt.arg, result, tt.want)
			}
		})
	}
}

func BenchmarkSumRec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sumRec(generateRandomSlice())
	}
}

func TestSumParallel(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SumParallel(tt.arg, 2)
			if result != tt.want {
				t.Errorf("sumParallel(%v) = %d; want %d", tt.arg, result, tt.want)
			}
		})
	}
}

func BenchmarkSumParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SumParallel(generateRandomSlice(), 2)
	}
}

func TestSumHeuristic(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SumHeuristic(tt.arg)
			if result != tt.want {
				t.Errorf("sumHeuristic(%v) = %d; want %d", tt.arg, result, tt.want)
			}
		})
	}
}

func BenchmarkSumHeuristic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SumHeuristic(generateRandomSlice())
	}
}

func TestSumHeuristicInline(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SumHeuristicInline(tt.arg)
			if result != tt.want {
				t.Errorf("sumHeuristicInline(%v) = %d; want %d", tt.arg, result, tt.want)
			}
		})
	}
}

func BenchmarkSumHeuristicInline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SumHeuristicInline(generateRandomSlice())
	}
}

func generateRandomSlice() []int {
	size := rand.Intn(maxSliceSize)
	slice := make([]int, size)
	for i := range slice {
		slice[i] = rand.Intn(maxItemValue)
	}
	return slice
}
