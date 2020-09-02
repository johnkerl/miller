package lib

func (this *Mlrval) IsAbsent() bool {
	return this.mvtype == MT_ABSENT
}
