package mlrval

// TODO: comment about mvtype; deferrence; copying of deferrence.
func (mv *Mlrval) Copy() *Mlrval {
	other := *mv
	if mv.mvtype == MT_MAP {
		other.intf = mv.intf.(*Mlrmap).Copy()
	} else if mv.mvtype == MT_ARRAY {
		other.intf = CopyMlrvalArray(mv.intf.([]*Mlrval))
	}
	return &other
}
