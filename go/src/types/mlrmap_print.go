package types

import (
	"bytes"
	"fmt"
	"os"
)

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) Print() {
	mlrmap.Fprint(os.Stdout)
	os.Stdout.WriteString("\n")
}
func (mlrmap *Mlrmap) Fprint(file *os.File) {
	(*file).WriteString(mlrmap.ToDKVPString())
}

func (mlrmap *Mlrmap) ToDKVPString() string {
	var buffer bytes.Buffer // stdio is non-buffered in Go, so buffer for ~5x speed increase
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(pe.Key)
		buffer.WriteString("=")
		buffer.WriteString(pe.Value.String())
		if pe.Next != nil {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}

// ----------------------------------------------------------------
// Must have non-pointer receiver in order to implement the fmt.Stringer
// interface to make mlrmap printable via fmt.Println et al.
func (mlrmap Mlrmap) String() string {
	bytes, err := mlrmap.MarshalJSON(JSON_MULTILINE, false)
	if err != nil {
		return "Mlrmap: could not not marshal self to JSON"
	} else {
		return string(bytes) + "\n"
	}
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) Dump() {
	fmt.Printf("FIELD COUNT: %d\n", mlrmap.FieldCount)
	fmt.Printf("HEAD:        %p\n", mlrmap.Head)
	fmt.Printf("TAIL:        %p\n", mlrmap.Tail)
	fmt.Printf("KEYS TO ENTRIES:\n")
	for k, e := range mlrmap.keysToEntries {
		fmt.Printf("  %-10s %#v\n", k, e)
	}
	fmt.Printf("LIST:\n")
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		fmt.Printf("  key: \"%s\" value: %-20s prev:%p self:%p next:%p\n",
			pe.Key, pe.Value.String(), pe.Prev, pe, pe.Next,
		)
	}
}
