package mlrval

// TODO: comment about mvtype; deferrence; copying of deferrence.
func (mv *Mlrval) Copy() *Mlrval {
	other := *mv
	if mv.mvtype == MT_MAP {
		other.mapval = mv.mapval.Copy()
	} else if mv.mvtype == MT_ARRAY {
		other.arrayval = CopyMlrvalArray(mv.arrayval)
	}
	return &other
}
