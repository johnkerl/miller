package main

import (
	// System:
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	// Miller:
	"miller/stream"
	// Temp:
	"miller/dsl"
	"miller/parsing/lexer"
	"miller/parsing/parser"
)

// ----------------------------------------------------------------
// xxx to do: stdout/stderr w/ ternary on exitrc
func usage() {
	// xxx temp grammar
	fmt.Fprintf(os.Stderr, "Usage: %s [options] {ifmt} {mapper} {ofmt} {filenames ...}\n",
		os.Args[0])
	fmt.Fprintf(os.Stderr, "If no file name is given, or if filename is \"-\", stdin is used.\n")
	flag.PrintDefaults()
	os.Exit(1)
}

// ----------------------------------------------------------------
func main() {
	runtime.GOMAXPROCS(4) // Seems reasonable these days

	cpuprofile := flag.String("cpuprofile", "", "Write CPU profile to `file`")

	flag.Usage = usage
	flag.Parse()
	maybeProfile(cpuprofile)

	args := flag.Args()

	// xxx temp
	if len(args) == 2 && args[0] == "parse" {
		parse(args[1])
		return
	}

	if len(args) < 3 {
		usage()
	}
	inputFormatName := args[0]
	mapperName := args[1]
	outputFormatName := args[2]
	filenames := args[3:]

	err := stream.Stream(filenames, inputFormatName, mapperName, outputFormatName)
	if err != nil {
		log.Fatal(err)
	}
}

// ----------------------------------------------------------------
// xxx temp
func parse(sourceString string) {
	fmt.Printf("Input: %s\n", sourceString)
	theLexer := lexer.NewLexer([]byte(sourceString))
	theParser := parser.NewParser()
	interfaceAST, err := theParser.Parse(theLexer)
	if err == nil {
		interfaceAST.(*dsl.AST).Print()
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ----------------------------------------------------------------
func maybeProfile(cpuprofile *string) {
	// to do: move to method
	// go tool pprof mlr foo.prof
	//   top10
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

}
