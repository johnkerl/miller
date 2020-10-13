package lib

import (
	"fmt"
	"math"
	"time"
)

func Sec2GMT(epochSeconds float64, numDecimalPlaces int) string {
	if numDecimalPlaces > 9 {
		numDecimalPlaces = 9
	}

	intPart := int64(epochSeconds)
	fractionalPart := epochSeconds - float64(intPart)
	decimalPart := int64(fractionalPart * math.Pow(10.0, float64(numDecimalPlaces)))
	t := time.Unix(intPart, 0).UTC()

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
