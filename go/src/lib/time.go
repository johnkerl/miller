package lib

import (
	"fmt"
	"math"
	"os"
	"time"
)

// SetTZFromEnv applies the $TZ environment variable. This has three reasons:
// (1) On Windows (as of 2021-10-20), this is necessary to get $TZ into use.
// (2) On Linux/Mac, as of this writing it is not necessary for initial value
// of TZ at startup. However, an explicit check is helpful since if someone
// does 'export TZ=Something/Invalid', then runs Miller, and invalid TZ is
// simply *ignored* -- we want to surface that error to the user.  (3) On any
// platform this is necessary for *changing* TZ mid-process: e.g.  if a DSL
// statement does 'ENV["TZ"] = Asia/Istanbul'.
func SetTZFromEnv() {
	location, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
		os.Exit(1)
	}
	time.Local = location
}

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
	if fractionalPart < 0 {
		intPart -= 1
		fractionalPart += 1.0
	}
	decimalPart := int64(fractionalPart * math.Pow(10.0, float64(numDecimalPlaces)))

	t := time.Unix(intPart, 0)
	if doLocal {
		// Note: the Go time package doesn't do Getenv("TZ") on every call.
		// Rather, it stashes the first call. This means we can't change $TZ
		// mid-process for testing purposes.
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
		if doLocal {
			return fmt.Sprintf(
				"%04d-%02d-%02d %02d:%02d:%02d",
				YYYY, MM, DD, hh, mm, ss)
		} else {
			return fmt.Sprintf(
				"%04d-%02d-%02dT%02d:%02d:%02dZ",
				YYYY, MM, DD, hh, mm, ss)
		}
	} else {
		if doLocal {
			return fmt.Sprintf(
				"%04d-%02d-%02d %02d:%02d:%02d.%0*d",
				YYYY, MM, DD, hh, mm, ss, numDecimalPlaces, decimalPart)
		} else {
			return fmt.Sprintf(
				"%04d-%02d-%02dT%02d:%02d:%02d.%0*dZ",
				YYYY, MM, DD, hh, mm, ss, numDecimalPlaces, decimalPart)
		}
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
