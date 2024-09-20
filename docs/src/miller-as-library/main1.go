package main

import (
	"fmt"

	"github.com/johnkerl/miller/v6/pkg/bifs"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

func main() {
	a := mlrval.FromInt(2)
	b := mlrval.FromInt(60)
	c := bifs.BIF_pow(a, b)
	fmt.Println(c.String())
}
