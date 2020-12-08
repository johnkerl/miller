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
	bytes, err := this.MarshalJSON()
	if err != nil {
		return "Mlrmap: could not not marshal self to JSON"
	} else {
		return string(bytes) + "\n"
	}
}

//// ----------------------------------------------------------------
//void lrec_dump(Mlrmap* prec) {
//	lrec_dump_fp(prec, stdout);
//}

//void lrec_dump_fp(Mlrmap* prec, FILE* fp) {
//	if (prec == NULL) {
//		fprintf(fp, "NULL\n");
//		return;
//	}
//	fprintf(fp, "field_count = %d\n", prec->field_count);
//	fprintf(fp, "| Head: %16p | Tail %16p\n", prec->Head, prec->Tail);
//	for (mlrmapEntry* pe = prec->Head; pe != NULL; pe = pe->Next) {
//		const char* key_string = (pe == NULL) ? "none" :
//			pe->key == NULL ? "null" :
//			pe->key;
//		const char* value_string = (pe == NULL) ? "none" :
//			pe->value == NULL ? "null" :
//			pe->value;
//		fprintf(fp,
//		"| prev: %16p curr: %16p next: %16p | key: %12s | value: %12s |\n",
//			pe->Prev, pe, pe->Next,
//			key_string, value_string);
//	}
//}

//void lrec_dump_titled(char* msg, Mlrmap* prec) {
//	printf("%s:\n", msg);
//	lrec_dump(prec);
//	printf("\n");
//}

//void lrec_print(Mlrmap* prec) {
//	FILE* output_stream = stdout;
//	char ors = '\n';
//	char ofs = ',';
//	char ops = '=';
//	if (prec == NULL) {
//		fputs("NULL", output_stream);
//		fputc(ors, output_stream);
//		return;
//	}
//	int nf = 0;
//	for (mlrmapEntry* pe = prec->Head; pe != NULL; pe = pe->Next) {
//		if (nf > 0)
//			fputc(ofs, output_stream);
//		fputs(pe->key, output_stream);
//		fputc(ops, output_stream);
//		fputs(pe->value, output_stream);
//		nf++;
//	}
//	fputc(ors, output_stream);
//}

//char* lrec_sprint(Mlrmap* prec, char* ors, char* ofs, char* ops) {
//	string_builder_t* psb = sb_alloc(SB_ALLOC_LENGTH);
//	if (prec == NULL) {
//		sb_append_string(psb, "NULL");
//	} else {
//		int nf = 0;
//		for (mlrmapEntry* pe = prec->Head; pe != NULL; pe = pe->Next) {
//			if (nf > 0)
//				sb_append_string(psb, ofs);
//			sb_append_string(psb, pe->key);
//			sb_append_string(psb, ops);
//			sb_append_string(psb, pe->value);
//			nf++;
//		}
//		sb_append_string(psb, ors);
//	}
//	char* rv = sb_finish(psb);
//	sb_free(psb);
//	return rv;
//}
