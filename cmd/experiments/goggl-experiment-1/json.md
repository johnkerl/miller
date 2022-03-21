# JSON Parser Example

This example is based on [the JSON spec](https://www.json.org/). It's intended for illustrative purposes, it hasn't been extensively tested.

To run this example, run `gogll json.md` and then `go test` on this folder. 

## GoGLL Header
```
package "github.com/johnkerl/miller/cmd/experiments/goggl-experiment-1"

GoGLL: Value;
``` 
## Lexical Rules
```
string : '"' 
            { not "␀␁␂␃␄␅␆␇␈␉␊␋␌␍␎␏␐␑␒␓␔␕␖␗␘␙␚␛␜␝␞␟\"\\"
            | '\\' any "\"\\/bfnrt" 
            | '\\' 'u' (number|'A'|'B'|'C'|'D'|'E'|'F'|'a'|'b'|'c'|'d'|'e'|'f') 
            } 
         '"';

numeric: ['-'] ( '0' ['.' <number>] | (any "123456789") {number} ['.' < number >])  
               [('E' | 'e')('+' | '-') <number>] ;
```

## Syntax Rules

```
Array : "[" "]" | "[" Values "]" ;
Values: Value | Value "," Values ;
Value : string | numeric | Object | Array | "true" | "false" | "null" ;

Object: "{" "}" | "{" Members "}" ;
Members: Member | Member "," Members ;
Member: string ":" Value ;
```
