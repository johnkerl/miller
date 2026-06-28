package mlrval

// TODO: comment about mvtype; deferrence; copying of deferrence.
func (mv *Mlrval) Copy() *Mlrval {
	other := *mv
	switch mv.mvtype {
	case MT_MAP:
		other.intf = mv.intf.(*Mlrmap).Copy()
	case MT_ARRAY:
		other.intf = CopyMlrvalArray(mv.intf.([]*Mlrval))
	}
	return &other
}
