package mlrval

import "bytes"

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

func (mlrmap *Mlrmap) ToNIDXString() string {
	var buffer bytes.Buffer // stdio is non-buffered in Go, so buffer for ~5x speed increase
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(pe.Value.String())
		if pe.Next != nil {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}

// Must have non-pointer receiver in order to implement the fmt.Stringer
// interface to make mlrmap printable via fmt.Println et al.
func (mlrmap Mlrmap) String() string {
	bytes, err := mlrmap.MarshalJSON(JSON_MULTILINE, false)
	if err != nil {
		return "Mlrmap: could not not marshal self to JSON"
	}
	return string(bytes) + "\n"
}
