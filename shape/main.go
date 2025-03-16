package main

import "fmt"

type Shape interface {
	Area() int
}

type Rectangle struct {
	Width  int
	Height int
}

func (r Rectangle) Area() int {
	return r.Height * r.Width
}

func (r *Rectangle) Scale(factor int) {
	r.Width *= factor
	r.Height *= factor
}

type Circle struct {
	Radius int
}

func (c Circle) Area() int {
	return c.Radius * c.Radius * 3
}

func (c *Circle) Scale(factor int) {
	c.Radius *= factor
}

type Square struct {
	Rectangle
}

func PrintArea(shape Shape) {
	fmt.Println(shape.Area())
}

func main() {
	rect := Rectangle{Width: 10, Height: 5}
	PrintArea(rect)
	rect.Scale(2)
	PrintArea(rect)

	circle := Circle{Radius: 7}
	PrintArea(circle)
	circle.Scale(2)
	PrintArea(circle)

	square := Square{Rectangle: Rectangle{Width: 4, Height: 4}}
	PrintArea(square)
	square.Scale(3)
	PrintArea(square)
	fmt.Println(square.Width, square.Height)
}
