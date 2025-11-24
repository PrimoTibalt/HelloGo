package main

import (
	"fmt"
	"math"
)

func main() {
	printTwoNumbers()
	constantAssignment()
	overflowCheck()
}

func printTwoNumbers() {
	i := 20
	f := float64(i)
	fmt.Println(i)
	fmt.Println(f)
}

func constantAssignment() {
	const number = 10
	var f float64 = number
	var i int = number
	fmt.Println(i)
	fmt.Println(f)
}

func overflowCheck() {
	var b byte = math.MaxUint8
	var smallI int32 = math.MaxInt32
	var bigI uint64 = math.MaxUint64

	b += 1
	smallI += 1
	bigI += 1

	fmt.Println(b)
	fmt.Println(smallI)
	fmt.Println(bigI)
}
