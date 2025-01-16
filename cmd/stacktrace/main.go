package main

import (
	"fmt"
	"runtime/debug"
)

func main() {
	log()

	fmt.Println("---")
	fmt.Println("---")
	fmt.Println("---")
}

func log() {
	debug.PrintStack()
}
