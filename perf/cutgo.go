package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

// ----------------------------------------------------------------
func main() {
	includeFields := strings.Split(os.Args[1], ",")
	filenames := os.Args[2:]

	ok := true
	if len(filenames) == 0 {
		ok = handle("-", includeFields) && ok
	} else {
		for _, arg := range filenames {
			ok = handle(arg, includeFields) && ok
		}
	}
	if ok {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// ----------------------------------------------------------------
func handle(fileName string, includeFields []string) (ok bool) {
	inputStream := os.Stdin
	if fileName != "-" {
		var err error
		if inputStream, err = os.Open(fileName); err != nil {
			log.Println(err)
			return false
		}
	}

	reader := bufio.NewReader(inputStream)
	eof := false

	for !eof {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			log.Println(err)
			if fileName != "-" {
				inputStream.Close()
			}
			return false
		} else {

			// continue
			// $ repeat 10 justtime cutgo ccode,milex,year,cinc ../c/nmc1.dkvp > /dev/null
			// TIME IN SECONDS 0.228 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 0.226 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 0.222 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 0.221 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 0.228 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 0.226 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 0.225 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 0.227 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 0.222 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 0.223 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			//
			// avg 0.2248
			// cumu  4%

			// Line to map
			mymap := make(map[string]string)
			fields := strings.Split(line, ",")
			for _, field := range fields {
				kvps := strings.SplitN(field, "=", 2)
				mymap[kvps[0]] = kvps[1]
			}

			// continue
			// $ repeat 10 justtime cutgo ccode,milex,year,cinc ../c/nmc1.dkvp > /dev/null
			// TIME IN SECONDS 3.055 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.837 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.905 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.817 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.766 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.810 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.748 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.744 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.765 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.722 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			//
			// avg 2.8169
			// cumu 58%
			// delta 54%

			// Map-to-map transform
			newmap := make(map[string]string)
			for _, includeField := range includeFields {
				value, present := mymap[includeField]
				if present {
					newmap[includeField] = value
				}
			}

			// continue
			// $ repeat 10 justtime cutgo ccode,milex,year,cinc ../c/nmc1.dkvp > /dev/null
			// TIME IN SECONDS 3.101 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.028 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.992 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.005 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.991 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.986 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.992 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.992 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.989 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 2.991 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			//
			// avg 3.0067
			// cumu 62%
			// delta  4%

			// Map to string

			// Faster to assemble in memory with single fmt.Println at the end,
			// than multiple fmt.Print through the fields.
			var buffer bytes.Buffer
			i := 0
			for k, v := range newmap {
				if i > 0 {
					buffer.WriteString(",")
				}
				buffer.WriteString(k)
				buffer.WriteString("=")
				buffer.WriteString(v)
				i++
			}

			// continue
			// $ repeat 10 justtime cutgo ccode,milex,year,cinc ../c/nmc1.dkvp > /dev/null
			// TIME IN SECONDS 3.821 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.443 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.491 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.421 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.412 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.470 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.438 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.555 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.637 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.538 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			//
			// avg 3.5226
			// cumu 73%
			// delta 11%

			buffer.WriteString("\n")

			// continue
			// $ repeat 10 justtime cutgo ccode,milex,year,cinc ../c/nmc1.dkvp > /dev/null
			// TIME IN SECONDS 3.769 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.539 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.539 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.486 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.518 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.492 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.506 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.486 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.495 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 3.560 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			//
			// avg 3.539
			// cumu 73%
			// delta  0%

			os.Stdout.WriteString(buffer.String())

			// $ repeat 10 justtime cutgo ccode,milex,year,cinc ../c/nmc1.dkvp > /dev/null
			// TIME IN SECONDS 4.976 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 4.727 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 4.760 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 4.747 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 4.826 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 4.849 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 4.841 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 4.747 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 4.841 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			// TIME IN SECONDS 4.765 -- cutgo ccode,milex,year,cinc ../c/nmc1.dkvp
			//
			// avg 4.8079
			// cumu 100%
			// delta 27%

		}
	}
	if fileName != "-" {
		inputStream.Close()
	}

	return true
}
