package hsd

import (
	"testing"
)

func TestStringDistance(t *testing.T) {
	tests := []struct {
		name string
		lhs  string
		rhs  string
		want int
	}{
		{name: "len(lhs) > len(rhs)", lhs: "foo", rhs: "fo", want: -1},
		{name: "len(lhs) < len(rhs)", lhs: "fo", rhs: "foo", want: -1},
		{name: "no diffs", lhs: "foo", rhs: "foo", want: 0},
		{name: "1 diffs", lhs: "foo", rhs: "foh", want: 1},
		{name: "2 diffs", lhs: "foo", rhs: "fhh", want: 2},
		{name: "all diffs", lhs: "foo", rhs: "bar", want: 3},
		{name: "empty", lhs: "", rhs: "", want: 0},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := StringDistance(test.lhs, test.rhs)
			if got != test.want {
				t.Fatalf("%s: expected: %v, got %v:", test.name, test.want, got)
			}
		})
	}
}
