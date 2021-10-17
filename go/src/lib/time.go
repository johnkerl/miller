package lib

import (
	"fmt"
	"math"
	"time"
)

func Sec2GMT(epochSeconds float64, numDecimalPlaces int) string {
	return sec2GMTOrLocalTime(epochSeconds, numDecimalPlaces, false)
}
func Sec2LocalTime(epochSeconds float64, numDecimalPlaces int) string {
	return sec2GMTOrLocalTime(epochSeconds, numDecimalPlaces, true)
}

func sec2GMTOrLocalTime(epochSeconds float64, numDecimalPlaces int, doLocal bool) string {
	if numDecimalPlaces > 9 {
		numDecimalPlaces = 9
	}

	intPart := int64(epochSeconds)
	fractionalPart := epochSeconds - float64(intPart)
	decimalPart := int64(fractionalPart * math.Pow(10.0, float64(numDecimalPlaces)))
	t := time.Unix(intPart, 0)
	if doLocal {
		t = t.Local()
	} else {
		t = t.UTC()
	}

	YYYY := t.Year()
	MM := int(t.Month())
	DD := t.Day()
	hh := t.Hour()
	mm := t.Minute()
	ss := t.Second()

	if numDecimalPlaces == 0 {
		return fmt.Sprintf(
			"%04d-%02d-%02dT%02d:%02d:%02dZ",
			YYYY, MM, DD, hh, mm, ss)
	} else {
		return fmt.Sprintf(
			"%04d-%02d-%02dT%02d:%02d:%02d.%0*dZ",
			YYYY, MM, DD, hh, mm, ss, numDecimalPlaces, decimalPart)
	}
}

func EpochSecondsToGMT(epochSeconds float64) time.Time {
	return epochSecondsToGMTOrLocalTime(epochSeconds, false)
}

func EpochSecondsToLocalTime(epochSeconds float64) time.Time {
	return epochSecondsToGMTOrLocalTime(epochSeconds, true)
}

func epochSecondsToGMTOrLocalTime(epochSeconds float64, doLocal bool) time.Time {
	intPart := int64(epochSeconds)
	fractionalPart := epochSeconds - float64(intPart)
	decimalPart := int64(fractionalPart * 1e9)
	if doLocal {
		return time.Unix(intPart, decimalPart).Local()
	} else {
		return time.Unix(intPart, decimalPart).UTC()
	}
}
