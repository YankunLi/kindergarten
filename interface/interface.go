package main

import "fmt"

type Face interface {
	Show()
}

type RedFace struct {
	color string
}
type BlackFace struct {
	color string
}

func (r *RedFace) Show() {
	fmt.Println("This face is Red")
}

func (b *BlackFace) Show() {
	fmt.Println("This face is Black")
}

func Smile(face Face) {
	face.Show()
}

func GetFace(b bool) Face {
	if b {
		return &RedFace{color: "red"}
	} else {
		return &BlackFace{color: "black"}
	}
}

func describe(t Face) {
	fmt.Printf("Interface type : %T, value: %v\n", t, t)
}

func assert(t Face) {
	//	v := t.(*RedFace)
	//fmt.Println(v)
	v, ok := t.(*RedFace)
	fmt.Println(v, ok)
}

func findType(i Face) {
	switch i.(type) {
	case *RedFace:
		fmt.Printf("RedFace: %s\n", i.(*RedFace))
	case *BlackFace:
		fmt.Printf("BlackFace: %s\n", i.(*BlackFace))
	default:
		fmt.Printf("Unknown type\n")
	}
}
func main() {
	var ff = []Face{&RedFace{}, &BlackFace{}}

	for _, val := range ff {
		Smile(val)
	}

	face := GetFace(true)
	face.Show()
	describe(face)

	assert(&RedFace{color: "Red"})
	assert(&BlackFace{color: "Black"})

	findType(face)

}
