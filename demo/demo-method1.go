package demo

import "fmt"

type Book struct {
	pages int
}

func (b *Book) SetPages(pages int) {
	b.pages = pages
}

func (a Book) SetPages1(pages int){
	a.pages = pages
}

func main() {
	var b *Book = new(Book)
	var a Book
	b.SetPages(1)
	fmt.Println(b.pages)
	a.SetPages1(1)
	fmt.Println(a.pages)
}