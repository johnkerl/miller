package types

import (
	"bytes"
	"fmt"
	"os"
)

// ----------------------------------------------------------------
func (this *Mlrmap) Print() {
	this.Fprint(os.Stdout)
	os.Stdout.WriteString("\n")
}
func (this *Mlrmap) Fprint(file *os.File) {
	(*file).WriteString(this.ToDKVPString())
}

func (this *Mlrmap) ToDKVPString() string {
	var buffer bytes.Buffer // 5x faster than fmt.Print() separately
	for pe := this.Head; pe != nil; pe = pe.Next {
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
// interface to make this printable via fmt.Println et al.
func (this Mlrmap) String() string {
	bytes, err := this.MarshalJSON(JSON_MULTILINE)
	if err != nil {
		return "Mlrmap: could not not marshal self to JSON"
	} else {
		return string(bytes) + "\n"
	}
}

// ----------------------------------------------------------------
func (this *Mlrmap) Dump() {
	fmt.Printf("FIELD COUNT: %d\n", this.FieldCount)
	fmt.Printf("HEAD:        %p\n", this.Head)
	fmt.Printf("TAIL:        %p\n", this.Tail)
	fmt.Printf("KEYS TO ENTRIES:\n")
	for k, e := range this.keysToEntries {
		fmt.Printf("  %-10s %#v\n", k, e)
	}
	fmt.Printf("LIST:\n")
	for pe := this.Head; pe != nil; pe = pe.Next {
		fmt.Printf("  key: \"%s\" value: %-20s prev:%p self:%p next:%p\n",
			pe.Key, pe.Value.String(), pe.Prev, pe, pe.Next,
		)
	}
}
