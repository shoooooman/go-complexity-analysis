package a

import (
	"fmt"
)

func comp1() { // want "Cyclomatic complexity: 1"
	var a int
	a++
	print(a)
}

func comp2() { // want "Cyclomatic complexity: 1"
	defer fmt.Println("world")

	fmt.Println("hello")
}

func comp3() { // want "Cyclomatic complexity: 1"
	go fmt.Println("hello")
	fmt.Println("world")
}

func comp4() { // want "Cyclomatic complexity: 1"
	a := make(chan string)
	go func() { a <- "ping" }()

	b := <-a
	fmt.Println(b)
}

func comp5() { // want "Cyclomatic complexity: 1"
	fmt.Println("Hello, 世界!")
	return
}

func comp6() { // want "Cyclomatic complexity: 3"
	var a int
	for a < 5 {
		if a < 3 {
			continue
		} else {
			break
		}
		a++
	}
}

func comp7() { // want "Cyclomatic complexity: 3"
	c1 := make(chan string)

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		}
	}
}
func comp8() { // want "Cyclomatic complexity: 2"
	a := []int{0, 1, 2}
	for b := range a {
		fmt.Println(b)
	}
}
