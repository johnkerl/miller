package bifs

import (
	"fmt"
	"math"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

func BIF_dhms2sec(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsString() {
		return mlrval.ERROR
	}
	var d, h, m, s int64

	if strings.HasPrefix(input1.AcquireStringValue(), "-") {

		n, err := fmt.Sscanf(input1.AcquireStringValue(), "-%dd%dh%dm%ds", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return mlrval.FromInt(-(s + m*60 + h*60*60 + d*60*60*24))
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "-%dh%dm%ds", &h, &m, &s)
		if n == 3 && err == nil {
			return mlrval.FromInt(-(s + m*60 + h*60*60))
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "-%dm%ds", &m, &s)
		if n == 2 && err == nil {
			return mlrval.FromInt(-(s + m*60))
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "-%ds", &s)
		if n == 1 && err == nil {
			return mlrval.FromInt(-(s))
		}

	} else {

		n, err := fmt.Sscanf(input1.AcquireStringValue(), "%dd%dh%dm%ds", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return mlrval.FromInt(s + m*60 + h*60*60 + d*60*60*24)
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "%dh%dm%ds", &h, &m, &s)
		if n == 3 && err == nil {
			return mlrval.FromInt(s + m*60 + h*60*60)
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "%dm%ds", &m, &s)
		if n == 2 && err == nil {
			return mlrval.FromInt(s + m*60)
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "%ds", &s)
		if n == 1 && err == nil {
			return mlrval.FromInt(s)
		}

	}
	return mlrval.ERROR
}

func BIF_dhms2fsec(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsString() {
		return mlrval.ERROR
	}

	var d, h, m int
	var s float64

	if strings.HasPrefix(input1.AcquireStringValue(), "-") {

		n, err := fmt.Sscanf(input1.AcquireStringValue(), "-%dd%dh%dm%fs", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return mlrval.FromFloat(-(s + float64(m*60+h*60*60+d*60*60*24)))
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "-%dh%dm%fs", &h, &m, &s)
		if n == 3 && err == nil {
			return mlrval.FromFloat(-(s + float64(m*60+h*60*60)))
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "-%dm%fs", &m, &s)
		if n == 2 && err == nil {
			return mlrval.FromFloat(-(s + float64(m*60)))
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "-%fs", &s)
		if n == 1 && err == nil {
			return mlrval.FromFloat(-(s))
		}

	} else {

		n, err := fmt.Sscanf(input1.AcquireStringValue(), "%dd%dh%dm%fs", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return mlrval.FromFloat(s + float64(m*60+h*60*60+d*60*60*24))
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "%dh%dm%fs", &h, &m, &s)
		if n == 3 && err == nil {
			return mlrval.FromFloat(s + float64(m*60+h*60*60))
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "%dm%fs", &m, &s)
		if n == 2 && err == nil {
			return mlrval.FromFloat(s + float64(m*60))
		}
		n, err = fmt.Sscanf(input1.AcquireStringValue(), "%fs", &s)
		if n == 1 && err == nil {
			return mlrval.FromFloat(s)
		}

	}
	return mlrval.ERROR
}

func BIF_hms2sec(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsString() {
		return mlrval.ERROR
	}
	if input1.AcquireStringValue() == "" {
		return mlrval.ERROR
	}
	var h, m, s int64

	if strings.HasPrefix(input1.AcquireStringValue(), "-") {
		n, err := fmt.Sscanf(input1.AcquireStringValue(), "-%d:%d:%d", &h, &m, &s)
		if n == 3 && err == nil {
			return mlrval.FromInt(-(s + m*60 + h*60*60))
		}
	} else {
		n, err := fmt.Sscanf(input1.AcquireStringValue(), "%d:%d:%d", &h, &m, &s)
		if n == 3 && err == nil {
			return mlrval.FromInt(s + m*60 + h*60*60)
		}
	}

	return mlrval.ERROR
}

func BIF_hms2fsec(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsString() {
		return mlrval.ERROR
	}

	var h, m int
	var s float64

	if strings.HasPrefix(input1.AcquireStringValue(), "-") {
		n, err := fmt.Sscanf(input1.AcquireStringValue(), "-%d:%d:%f", &h, &m, &s)
		if n == 3 && err == nil {
			return mlrval.FromFloat(-(s + float64(m*60+h*60*60)))
		}
	} else {
		n, err := fmt.Sscanf(input1.AcquireStringValue(), "%d:%d:%f", &h, &m, &s)
		if n == 3 && err == nil {
			return mlrval.FromFloat(s + float64(m*60+h*60*60))
		}
	}

	return mlrval.ERROR
}

func BIF_sec2dhms(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	isec, ok := input1.GetIntValue()
	if !ok {
		return mlrval.ERROR
	}

	var d, h, m, s int64

	splitIntToDHMS(isec, &d, &h, &m, &s)
	if d != 0 {
		return mlrval.FromString(
			fmt.Sprintf("%dd%02dh%02dm%02ds", d, h, m, s),
		)
	} else if h != 0 {
		return mlrval.FromString(
			fmt.Sprintf("%dh%02dm%02ds", h, m, s),
		)
	} else if m != 0 {
		return mlrval.FromString(
			fmt.Sprintf("%dm%02ds", m, s),
		)
	} else {
		return mlrval.FromString(
			fmt.Sprintf("%ds", s),
		)
	}
}

func BIF_sec2hms(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	isec, ok := input1.GetIntValue()
	if !ok {
		return mlrval.ERROR
	}
	sign := ""
	if isec < 0 {
		sign = "-"
		isec = -isec
	}

	var d, h, m, s int64

	splitIntToDHMS(isec, &d, &h, &m, &s)
	h += d * 24

	return mlrval.FromString(
		fmt.Sprintf("%s%02d:%02d:%02d", sign, h, m, s),
	)
}

func BIF_fsec2dhms(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	fsec, ok := input1.GetNumericToFloatValue()
	if !ok {
		return mlrval.ERROR
	}

	sign := int64(1)
	if fsec < 0 {
		sign = -1
		fsec = -fsec
	}
	isec := int64(math.Trunc(fsec))
	fractional := fsec - float64(isec)

	var d, h, m, s int64

	splitIntToDHMS(isec, &d, &h, &m, &s)

	if d != 0 {
		d = sign * d
		return mlrval.FromString(
			fmt.Sprintf(
				"%dd%02dh%02dm%09.6fs",
				d, h, m, float64(s)+fractional),
		)
	} else if h != 0 {
		h = sign * h
		return mlrval.FromString(
			fmt.Sprintf(
				"%dh%02dm%09.6fs",
				h, m, float64(s)+fractional),
		)
	} else if m != 0 {
		m = sign * m
		return mlrval.FromString(
			fmt.Sprintf(
				"%dm%09.6fs",
				m, float64(s)+fractional),
		)
	} else {
		s = sign * s
		fractional = float64(sign) * fractional
		return mlrval.FromString(
			fmt.Sprintf(
				"%.6fs",
				float64(s)+fractional),
		)
	}
}

func BIF_fsec2hms(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	fsec, ok := input1.GetNumericToFloatValue()
	if !ok {
		return mlrval.ERROR
	}

	sign := ""
	if fsec < 0 {
		sign = "-"
		fsec = -fsec
	}
	isec := int64(math.Trunc(fsec))
	fractional := fsec - float64(isec)

	var d, h, m, s int64

	splitIntToDHMS(isec, &d, &h, &m, &s)
	h += d * 24

	// "%02.6f" does not exist so we have to do our own zero-pad
	if s < 10 {
		return mlrval.FromString(
			fmt.Sprintf("%s%02d:%02d:0%.6f", sign, h, m, float64(s)+fractional),
		)
	} else {
		return mlrval.FromString(
			fmt.Sprintf("%s%02d:%02d:%.6f", sign, h, m, float64(s)+fractional),
		)
	}
}

// Helper function
func splitIntToDHMS(u int64, pd, ph, pm, ps *int64) {
	d := int64(0)
	h := int64(0)
	m := int64(0)
	s := int64(0)
	sign := int64(1)
	if u < 0 {
		u = -u
		sign = -1
	}
	s = u % 60
	u = u / 60
	if u == 0 {
		s = s * sign
	} else {
		m = u % 60
		u = u / 60
		if u == 0 {
			m = m * sign
		} else {
			h = u % 24
			u = u / 24
			if u == 0 {
				h = h * sign
			} else {
				d = u * sign
			}
		}
	}
	*pd = d
	*ph = h
	*pm = m
	*ps = s
}
