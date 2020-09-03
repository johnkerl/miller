package lib

func (this *Mlrval) IsAbsent() bool {
	return this.mvtype == MT_ABSENT
}

func (this *Mlrval) GetBoolValue() (boolValue bool, isBoolean bool) {
	if this.mvtype == MT_BOOL {
		return this.boolval, true
	} else {
		return false, false
	}
}
