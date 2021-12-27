package main

import (
	"fmt"
)

var levelstep = 5
var chars = "X*o-."

func main() {
	rcorn := -2.0
	icorn := -2.0
	side := 4.0
	iheight := 500
	iwidth := 1000
	maxits := 100
	silent := false

	for ii := iheight - 1; ii >= 0; ii-- {
		ci := icorn + (float64(ii)/float64(iheight))*side
		for ir := 0; ir < iwidth; ir++ {
			cr := rcorn + (float64(ir)/float64(iwidth))*side
			c := get_point_plot(cr, ci, maxits)
			if !silent {
				fmt.Print(c)
			}
		}
		if !silent {
			fmt.Println()
		}
	}
}

func get_point_plot(cr float64, ci float64, maxits int) string {
	zr := 0.0
	zi := 0.0

	escaped := false
	zt := 0.0
	iti := 0
	for iti = 0; iti < maxits; iti++ {
		mag := zr*zr + zi + zi
		if mag > 4.0 {
			escaped = true
			break
		}
		// z := z^2 + c
		zt = zr*zr - zi*zi + cr
		zi = 2*zr*zi + ci
		zr = zt
	}
	if !escaped {
		return "."
	} else {
		level := (iti / levelstep) % len(chars)
		return chars[level : level+1]
	}
}
