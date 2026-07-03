package mlrval

// TODO: comment about mvtype; deferrence; copying of deferrence.
func (mv *Mlrval) Copy() *Mlrval {
	other := *mv
	switch mv.mvtype {
	case MT_MAP:
		other.intf = mv.intf.(*Mlrmap).Copy()
	case MT_ARRAY:
		other.intf = CopyMlrvalArray(mv.intf.([]*Mlrval))
	case MT_BYTES:
		other.intf = append([]byte(nil), mv.intf.([]byte)...)
	}
	return &other
}
