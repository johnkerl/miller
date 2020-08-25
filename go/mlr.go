package main

import (
	//"fmt"
	"containers"
)

func main() {
	lrec := containers.LrecAlloc()
	lrec.Put("a", "foo")
	lrec.Put("x", "3")
	lrec.Put("y", "y")
	lrec.Print()
}
