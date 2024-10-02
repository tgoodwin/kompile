package main

import "fmt"

func sayHi() {
	fmt.Println("Hi!")
}

func runTest() {
	fmt.Println("Hello, World!")

	go sayHi()
}
