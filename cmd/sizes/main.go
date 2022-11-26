// ================================================================
// Experiments for type-inference performance optimization
// ================================================================

// go build github.com/johnkerl/miller/cmd/sizes

package main

import (
	"fmt"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

func main() {
	var mvs [2]mlrval.Mlrval
	mvs[0] = *mlrval.FromString("hello")
	mvs[1] = *mlrval.FromString("world")
	mvs[0].ShowSizes()
	fmt.Println()
	mvs[1].ShowSizes()
}
