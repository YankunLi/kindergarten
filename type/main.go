package main

import "fmt"

type Bag struct {
	Key  string
	Size int
}

type Bag1 struct {
	Key  int
	Size int
}

func main() {
	var b1 interface{}
	var b2 interface{}

	b1 = Bag{Key: "1"}
	b2 = Bag1{Key: 1}
	{
		b, ok := b1.(Bag)
		fmt.Println("Bag type: ", ok, " data: ", b)
	}
	{
		b, ok := b2.(Bag1)
		fmt.Println("Bag type: ", ok, " data: ", b)
	}
	{
		b, ok := b1.(Bag1)
		fmt.Println("Bag type: ", ok, " data: ", b)
	}
}
