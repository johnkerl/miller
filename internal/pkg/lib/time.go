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
	return sec2Time(epochSeconds, numDecimalPlaces, false, nil)
}
func Sec2LocalTime(epochSeconds float64, numDecimalPlaces int) string {
	return sec2Time(epochSeconds, numDecimalPlaces, true, nil)
}

func Sec2LocationTime(epochSeconds float64, numDecimalPlaces int, location *time.Location) string {
	return sec2Time(epochSeconds, numDecimalPlaces, true, location)
}

// sec2Time is for DSL functions sec2gmt and sec2localtime. If doLocal is
// false, use UTC.  Else if location is nil, use $TZ environment variable. Else
// use the specified location.
func sec2Time(epochSeconds float64, numDecimalPlaces int, doLocal bool, location *time.Location) string {
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
		if location != nil {
			t = t.In(location)
		} else {
			t = t.Local()
		}
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
	return epochSecondsToTime(epochSeconds, false, nil)
}

func EpochSecondsToLocalTime(epochSeconds float64) time.Time {
	return epochSecondsToTime(epochSeconds, true, nil)
}

func EpochSecondsToLocationTime(epochSeconds float64, location *time.Location) time.Time {
	return epochSecondsToTime(epochSeconds, true, location)
}

func epochSecondsToTime(epochSeconds float64, doLocal bool, location *time.Location) time.Time {
	intPart := int64(epochSeconds)
	fractionalPart := epochSeconds - float64(intPart)
	decimalPart := int64(fractionalPart * 1e9)
	if doLocal {
		if location == nil {
			return time.Unix(intPart, decimalPart).Local()
		} else {
			return time.Unix(intPart, decimalPart).In(location)
		}
	} else {
		return time.Unix(intPart, decimalPart).UTC()
	}
}
