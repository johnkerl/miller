package main

import (
	"fmt"

	"github.com/johnkerl/miller/cmd/experiments/goggl-experiment-1/lexer"
	"github.com/johnkerl/miller/cmd/experiments/goggl-experiment-1/parser"
)

var input = []rune(`{ "key1": "hello world", "key2": [-1234.5678E-9, 10] }`)

func main() {
	l := lexer.New(input)

	fmt.Println("---- Token list ----")

	for _, t := range l.Tokens {
		fmt.Println(t.LiteralString())
	}

	BSR, err := parser.Parse(l)
	if err != nil {
		for _, e := range err {
			fmt.Println(e)
		}
		return
	}

	fmt.Println("---- Syntax Tree ----")

	BSR.Dump()

	// Output:
	// ---- Token list ----
	// {
	// "key1"
	// :
	// "hello world"
	// ,
	// "key2"
	// :
	// [
	// -1234.5678E-9
	// ,
	// 10
	// ]
	// }
	//
	// ---- Syntax Tree ----
	//GoGLL : Value ∙,0,0,13 - { "key1": "hello world", "key2": [-1234.5678E-9, 10] }
	//     Value : Object ∙,0,0,13 - { "key1": "hello world", "key2": [-1234.5678E-9, 10] }
	//         Object : { Members } ∙,0,12,13 - { "key1": "hello world", "key2": [-1234.5678E-9, 10] }
	//             Members : Member , Members ∙,1,5,12 - "key1": "hello world", "key2": [-1234.5678E-9, 10]
	//                 Member : string : Value ∙,1,3,4 - "key1": "hello world"
	//                     Value : string ∙,3,3,4 - "hello world"
	//                 Members : Member ∙,5,5,12 - "key2": [-1234.5678E-9, 10]
	//                     Member : string : Value ∙,5,7,12 - "key2": [-1234.5678E-9, 10]
	//                         Value : Array ∙,7,7,12 - [-1234.5678E-9, 10]
	//                             Array : [ Values ] ∙,7,11,12 - [-1234.5678E-9, 10]
	//                                 Values : Value , Values ∙,8,10,11 - -1234.5678E-9, 10
	//                                     Value : numeric ∙,8,8,9 - -1234.5678E-9
	//                                     Values : Value ∙,10,10,11 - 10
	//                                         Value : numeric ∙,10,10,11 - 10
}
