package mlrval

// TODO: comment about mvtype; deferrence; copying of deferrence.
func (mv *Mlrval) Copy() *Mlrval {
	if mv.mvtype == MT_MAP {
		panic("mlrval map-valued copy unimplemented, pending refactor")
	} else if mv.mvtype == MT_ARRAY {
		panic("mlrval array-valued copy unimplemented, pending refactor")
	} else {
		other := *mv
		return &other
	}
}
