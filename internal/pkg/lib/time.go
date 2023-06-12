package lib

import (
	"fmt"
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
func SetTZFromEnv() error {
	tzenv := os.Getenv("TZ")
	location, err := time.LoadLocation(tzenv)
	if err != nil {
		return fmt.Errorf("TZ environment variable appears malformed: \"%s\"", tzenv)
	}
	time.Local = location
	return nil
}

func Sec2GMT(epochSeconds float64, numDecimalPlaces int) string {
	return secToFormattedTime(epochSeconds, numDecimalPlaces, false, nil)
}

func Nsec2GMT(epochNanoseconds int64, numDecimalPlaces int) string {
	return nsecToFormattedTime(epochNanoseconds, numDecimalPlaces, false, nil)
}

func Sec2LocalTime(epochSeconds float64, numDecimalPlaces int) string {
	return secToFormattedTime(epochSeconds, numDecimalPlaces, true, nil)
}

func Nsec2LocalTime(epochNanoseconds int64, numDecimalPlaces int) string {
	return nsecToFormattedTime(epochNanoseconds, numDecimalPlaces, true, nil)
}

func Sec2LocationTime(epochSeconds float64, numDecimalPlaces int, location *time.Location) string {
	return secToFormattedTime(epochSeconds, numDecimalPlaces, true, location)
}

func Nsec2LocationTime(epochNanoseconds int64, numDecimalPlaces int, location *time.Location) string {
	return nsecToFormattedTime(epochNanoseconds, numDecimalPlaces, true, location)
}

// secToFormattedTime is for DSL functions sec2gmt and sec2localtime. If doLocal is
// false, use UTC.  Else if location is nil, use $TZ environment variable. Else
// use the specified location.
func secToFormattedTime(epochSeconds float64, numDecimalPlaces int, doLocal bool, location *time.Location) string {
	intPart := int64(epochSeconds)
	fractionalPart := epochSeconds - float64(intPart)
	if fractionalPart < 0 {
		intPart -= 1
		fractionalPart += 1.0
	}

	t := time.Unix(intPart, int64(fractionalPart*1e9))
	return goTimeToFormattedTime(t, numDecimalPlaces, doLocal, location)
}

// nsecToFormattedTime is for DSL functions nsec2gmt and nsec2localtime. If doLocal is
// false, use UTC.  Else if location is nil, use $TZ environment variable. Else
// use the specified location.
func nsecToFormattedTime(epochNanoseconds int64, numDecimalPlaces int, doLocal bool, location *time.Location) string {
	t := time.Unix(epochNanoseconds/1000000000, epochNanoseconds%1000000000)
	return goTimeToFormattedTime(t, numDecimalPlaces, doLocal, location)
}

// This is how much to divide nanoseconds by to get a desired number of decimal places
var nsToFracDivisors = []int{
	/* 0 */ 0, /* unused */
	/* 1 */ 100000000,
	/* 2 */ 10000000,
	/* 3 */ 1000000,
	/* 4 */ 100000,
	/* 5 */ 10000,
	/* 6 */ 1000,
	/* 7 */ 100,
	/* 8 */ 10,
	/* 9 */ 1,
}

func goTimeToFormattedTime(t time.Time, numDecimalPlaces int, doLocal bool, location *time.Location) string {
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

	if numDecimalPlaces < 0 {
		numDecimalPlaces = 0
	} else if numDecimalPlaces > 9 {
		numDecimalPlaces = 9
	}

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
		fractionalPart := t.Nanosecond() / nsToFracDivisors[numDecimalPlaces]
		if doLocal {
			return fmt.Sprintf(
				"%04d-%02d-%02d %02d:%02d:%02d.%0*d",
				YYYY, MM, DD, hh, mm, ss, numDecimalPlaces, fractionalPart)
		} else {
			return fmt.Sprintf(
				"%04d-%02d-%02dT%02d:%02d:%02d.%0*dZ",
				YYYY, MM, DD, hh, mm, ss, numDecimalPlaces, fractionalPart)
		}
	}
}

func EpochSecondsToGMT(epochSeconds float64) time.Time {
	return epochSecondsToTime(epochSeconds, false, nil)
}

func EpochNanosecondsToGMT(epochNanoseconds int64) time.Time {
	return epochNanosecondsToTime(epochNanoseconds, false, nil)
}

func EpochSecondsToLocalTime(epochSeconds float64) time.Time {
	return epochSecondsToTime(epochSeconds, true, nil)
}

func EpochNanosecondsToLocalTime(epochNanoseconds int64) time.Time {
	return epochNanosecondsToTime(epochNanoseconds, true, nil)
}

func EpochSecondsToLocationTime(epochSeconds float64, location *time.Location) time.Time {
	return epochSecondsToTime(epochSeconds, true, location)
}

func EpochNanosecondsToLocationTime(epochNanoseconds int64, location *time.Location) time.Time {
	return epochNanosecondsToTime(epochNanoseconds, true, location)
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

func epochNanosecondsToTime(epochNanoseconds int64, doLocal bool, location *time.Location) time.Time {
	intPart := epochNanoseconds / 1000000000
	fractionalPart := epochNanoseconds % 1000000000
	if doLocal {
		if location == nil {
			return time.Unix(intPart, fractionalPart).Local()
		} else {
			return time.Unix(intPart, fractionalPart).In(location)
		}
	} else {
		return time.Unix(intPart, fractionalPart).UTC()
	}
}
