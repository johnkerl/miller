package types

import (
	"bytes"
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
		buffer.WriteString(*pe.Key)
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
