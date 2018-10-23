package main

import (
	"fmt"
	"sync"
)

type A struct {
	Name string
}

type B struct {
	Name string
	Age  int32
}

func main() {
	var pool = sync.Pool{
		New: func() interface{} { return A{Name: "hello world"} },
	}

	val := A{Name: "hello sync"}
	Bval := B{
		Name: "B sync",
		Age:  12,
	}

	pool.Put(val)
	pool.Put(Bval)

	fmt.Println(pool.Get())
	fmt.Println(pool.Get())
	fmt.Println(pool.Get())
}
