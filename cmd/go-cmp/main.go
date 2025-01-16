package main

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type Person struct {
	ID   int
	Name string
}

func main() {
	s1 := []Person{{1, "hoge"}, {2, "fuga"}, {3, "piyo"}}
	s2 := []Person{{1, "hoge"}, {2, "fuga"}, {3, "piyo"}}

	diff := cmp.Diff(s1, s2)
	fmt.Println("diff s1 s2", diff)

	s3 := []Person{{1, "hoge"}, {3, "piyo"}, {2, "fuga"}}
	diff = cmp.Diff(s1, s3, cmpopts.SortSlices(func(x, y Person) bool { return fmt.Sprintf("%v", x) < fmt.Sprintf("%v", y) }))
	fmt.Println("diff s1 s3", diff)
}
