package lib

import (
	"strconv"
)

// See mlrval.go for more about JIT-formatting of string backings
func (this *Mlrval) setPrintRep() {
	if !this.printrepValid {
		// xxx do it -- disposition vector
		// xxx temp temp temp temp temp
		switch this.mvtype {
		case MT_ERROR:
			this.printrep = "(error)" // xxx constdef at top of file
			break
		case MT_ABSENT:
			// Callsites should be using absence to do non-assigns, so flag
			// this clearly visually if it should (buggily) slip through to
			// user-level visibility.
			this.printrep = "(bug-if-you-see-this)" // xxx constdef at top of file
			break
		case MT_VOID:
			this.printrep = "" // xxx constdef at top of file
			break
		case MT_STRING:
			// panic i suppose
			break
		case MT_INT:
			this.printrep = strconv.FormatInt(this.intval, 10)
			break
		case MT_FLOAT:
			// xxx temp -- OFMT etc ...
			this.printrep = strconv.FormatFloat(this.floatval, 'g', -1, 64)
			break
		case MT_BOOL:
			if this.boolval == true {
				this.printrep = "true"
			} else {
				this.printrep = "false"
			}
			break
		}
		this.printrepValid = true
	}
}

// Must have non-pointer receiver in order to implement the fmt.Stringer
// interface to make this printable via fmt.Println et al.
func (this Mlrval) String() string {
	this.setPrintRep()
	return this.printrep
}

// For JSON output. Second return value is true if the mlrval should be
// double-quoted.
func (this *Mlrval) StringWithQuoteInfo() (string, bool) {
	this.setPrintRep()
	quoteless := (this.mvtype == MT_INT || this.mvtype == MT_FLOAT || this.mvtype == MT_BOOL)
	return this.printrep, !quoteless
}
