package gen

import (
	"github.com/goccmack/gocc/internal/util/gen/golang"
)

func Gen(outDir string) {
	golang.GenRune(outDir)
	golang.GenLitConv(outDir)
}
