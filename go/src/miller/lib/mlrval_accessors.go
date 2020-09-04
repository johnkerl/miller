package lib

func (this *Mlrval) GetType() MVType {
	return this.mvtype
}

func (this *Mlrval) IsError() bool {
	return this.mvtype == MT_ERROR
}

func (this *Mlrval) IsAbsent() bool {
	return this.mvtype == MT_ABSENT
}

func (this *Mlrval) IsVoid() bool {
	return this.mvtype == MT_VOID
}

func (this *Mlrval) IsErrorOrVoid() bool {
	return this.mvtype == MT_ERROR || this.mvtype == MT_VOID
}

func (this *Mlrval) IsBool() bool {
	return this.mvtype == MT_BOOL
}

func (this *Mlrval) GetBoolValue() (boolValue bool, isBoolean bool) {
	if this.mvtype == MT_BOOL {
		return this.boolval, true
	} else {
		return false, false
	}
}

func (this *Mlrval) IsTrue() bool {
	return this.mvtype == MT_BOOL && this.boolval == true
}
func (this *Mlrval) IsFalse() bool {
	return this.mvtype == MT_BOOL && this.boolval == false
}

func (this *Mlrval) IsArray() bool {
	return this.mvtype == MT_ARRAY
}

func (this *Mlrval) GetArrayIndex(index *Mlrval) Mlrval {
	if this.mvtype != MT_ARRAY {
		return MlrvalFromError()
	}
	if index.mvtype != MT_INT {
		return MlrvalFromError()
	}
	i := index.intval
	n := int64(len(this.arrayval))

	// TODO: document this (pythonic)
	if i < 0 && i > -n {
		i += n
	}

	if i < 0 || i >= n {
		return MlrvalFromError()
	}
	return this.arrayval[i]
}
