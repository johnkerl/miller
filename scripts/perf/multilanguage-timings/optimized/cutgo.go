package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const pipelineCap = 64

type lineJob struct {
	index int
	line  string
}

type outJob struct {
	index int
	out   string
}

// CLI: step (1-6) arg1, comma-separated field names arg2, filenames arg3+.
// Step controls how far the pipeline runs (for profiling): 1=read, 2=+parse, 3=+select, 4=+build, 5=+newline, 6=+write.
func main() {
	if len(os.Args) < 3 {
		log.Fatalf("usage: %s <step 1-6> <field1,field2,...> [file ...]\n", os.Args[0])
	}
	step, err := strconv.Atoi(os.Args[1])
	if err != nil || step < 1 || step > 6 {
		log.Fatalf("step must be 1-6, got %q", os.Args[1])
	}
	includeFields := strings.Split(os.Args[2], ",")
	filenames := os.Args[3:]
	if len(filenames) == 0 {
		filenames = []string{"-"}
	}

	ok := true
	for _, arg := range filenames {
		ok = handle(arg, step, includeFields) && ok
	}
	if ok {
		os.Exit(0)
	}
	os.Exit(1)
}

func handle(fileName string, step int, includeFields []string) (ok bool) {
	inputStream := os.Stdin
	if fileName != "-" {
		var err error
		if inputStream, err = os.Open(fileName); err != nil {
			log.Println(err)
			return false
		}
		defer inputStream.Close()
	}

	readCh := make(chan lineJob, pipelineCap)
	writeCh := make(chan outJob, pipelineCap)
	done := make(chan struct{})
	errCh := make(chan error, 1)

	const readBufSize = 256 * 1024
	reader := bufio.NewReaderSize(inputStream, readBufSize)

	// Reader goroutine
	go func() {
		index := 0
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				select {
				case errCh <- err:
				default:
				}
				close(readCh)
				return
			}
			if err == io.EOF && line == "" {
				close(readCh)
				return
			}
			readCh <- lineJob{index: index, line: line}
			index++
			if err == io.EOF {
				close(readCh)
				return
			}
		}
	}()

	// Processor goroutine
	go func() {
		mymap := make(map[string]string)
		newmap := make(map[string]string, len(includeFields))
		var buffer bytes.Buffer
		for job := range readCh {
			if step <= 1 {
				continue
			}
			// Step 2: line to map
			for k := range mymap {
				delete(mymap, k)
			}
			fields := strings.Split(job.line, ",")
			for _, field := range fields {
				kvps := strings.SplitN(field, "=", 2)
				if len(kvps) >= 2 {
					mymap[kvps[0]] = kvps[1]
				}
			}
			if step <= 2 {
				continue
			}
			// Step 3: map-to-map transform
			for k := range newmap {
				delete(newmap, k)
			}
			for _, includeField := range includeFields {
				if value, present := mymap[includeField]; present {
					newmap[includeField] = value
				}
			}
			if step <= 3 {
				continue
			}
			// Step 4–5: map to string + newline
			buffer.Reset()
			first := true
			for _, includeField := range includeFields {
				if value, present := newmap[includeField]; present {
					if !first {
						buffer.WriteString(",")
					}
					buffer.WriteString(includeField)
					buffer.WriteString("=")
					buffer.WriteString(value)
					first = false
				}
			}
			buffer.WriteString("\n")
			if step <= 5 {
				continue
			}
			// Step 6: send to writer
			writeCh <- outJob{index: job.index, out: buffer.String()}
		}
		close(writeCh)
	}()

	// Writer goroutine
	go func() {
		for job := range writeCh {
			if _, err := os.Stdout.WriteString(job.out); err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
		}
		close(done)
	}()

	select {
	case err := <-errCh:
		log.Println(err)
		return false
	case <-done:
		return true
	}
}

