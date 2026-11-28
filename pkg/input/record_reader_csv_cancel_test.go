package input

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/types"
)

func TestCSVReadersStopBeforeLaterFiles(t *testing.T) {
	for _, format := range []string{"csv", "csvlite", "tsv"} {
		for _, cancel := range []bool{false, true} {
			name := format + "/complete"
			if cancel {
				name = format + "/cancelled"
			}
			t.Run(name, func(t *testing.T) {
				dir := t.TempDir()
				first := filepath.Join(dir, "first.csv")
				second := filepath.Join(dir, "second.csv")
				separator := ","
				if format == "tsv" {
					separator = "\t"
				}
				data := "a" + separator + "b\n" + strings.Repeat("1"+separator+"2\n", 100)
				if err := os.WriteFile(first, []byte(data), 0600); err != nil {
					t.Fatal(err)
				}
				if !cancel {
					if err := os.WriteFile(second, []byte(data), 0600); err != nil {
						t.Fatal(err)
					}
				}
				options := cli.DefaultReaderOptions()
				options.InputFileFormat = format
				if err := cli.FinalizeReaderOptions(&options); err != nil {
					t.Fatal(err)
				}
				reader, err := Create(&options, 1)
				if err != nil {
					t.Fatal(err)
				}
				records := make(chan []*types.RecordAndContext)
				errors := make(chan error, 10)
				done := make(chan bool, 1)
				context := types.NewContext()
				timer := time.NewTimer(5 * time.Second)
				defer timer.Stop()
				go reader.Read([]string{first, second}, *context, records, errors, done)
				count := 0
				sentDone := false
				for {
					var batch []*types.RecordAndContext
					select {
					case batch = <-records:
					case <-timer.C:
						t.Fatal("reader did not send an end-of-stream marker")
					}
					for _, record := range batch {
						if record.EndOfStream {
							if cancel && count >= 100 {
								t.Errorf("cancelled reader consumed all %d records", count)
							}
							if !cancel && count != 200 {
								t.Errorf("read %d records, want 200", count)
							}
							select {
							case err := <-errors:
								t.Fatalf("unexpected later input error: %v", err)
							default:
							}
							return
						}
						count++
					}
					if cancel && !sentDone {
						done <- true
						sentDone = true
					}
				}
			})
		}
	}
}
