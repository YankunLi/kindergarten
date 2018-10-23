package main

import "fmt"

type Face interface {
	Show()
}

type ReadFace struct{}
type BlackFace struct{}

func (r ReadFace) Show() {
	fmt.Println("This face is Red")
}

func (b BlackFace) Show() {
	fmt.Println("This face is Black")
}

func Smile(face Face) {
	face.Show()
}

func main() {
	var ff = []Face{ReadFace{}, BlackFace{}}

	for _, val := range ff {
		Smile(val)
	}

}
